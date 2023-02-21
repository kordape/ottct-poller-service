package worker

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kordape/ottct-poller-service/internal/processor"
	"github.com/kordape/ottct-poller-service/pkg/logger"
)

const defaultTickInterval = 1 * time.Second

type Worker struct {
	tickInterval time.Duration
	log          logger.Interface

	running     int32
	stopChannel chan bool
}

type Option func(w *Worker)

func WithInterval(interval time.Duration) Option {
	return func(w *Worker) {
		w.tickInterval = interval
	}
}

func NewWorker(log logger.Interface, opts ...Option) (*Worker, error) {
	stopChan := make(chan bool)

	w := &Worker{
		tickInterval: defaultTickInterval,
		log:          log,
		stopChannel:  stopChan,
	}

	for _, opt := range opts {
		opt(w)
	}

	if err := w.validate(); err != nil {
		return w, fmt.Errorf("Worker validation: %v", err)
	}

	return w, nil
}

func (w *Worker) validate() error {
	if w.log == nil {
		return errors.New("log is nil")
	}

	return nil
}

func (w *Worker) Run() error {
	if w.Running() {
		// if worker is already running just exit
		return nil
	}

	if err := w.validate(); err != nil {
		return fmt.Errorf("Can't run worker. Validation error: %v", err)
	}

	atomic.StoreInt32(&w.running, 1)
	ticker := time.NewTicker(w.tickInterval)

	go func() {
		for {
			select {
			case <-w.stopChannel:
				ticker.Stop()
				w.log.Info("Stopping worker")
				return
			case <-ticker.C:
				// create processing task
				w.log.Info("Worker tick")
				w.process()
			}
		}
	}()

	return nil
}

func (w *Worker) Running() bool {
	return atomic.LoadInt32(&w.running) == 1
}

func (w *Worker) Stop() {
	defer func() {
		atomic.StoreInt32(&w.running, 0)
	}()

	w.stopChannel <- true
}

func (w *Worker) process() {
	endTime := time.Now()
	startTime := endTime.Add(-w.tickInterval)

	w.log.Info(fmt.Sprintf("Processing for interval: %s - %s", startTime.Format(time.UnixDate), endTime.Format(time.UnixDate)))

	// TODO: replace this line with fetched entities from db
	entities := []string{"foo", "bar"}

	resultsChannels := w.startFetching(context.Background(), entities)
	results := make(processor.JobResults, 0, len(entities))
	fetchingEnded := w.collectResults(&results, resultsChannels)

	select {
	case <-fetchingEnded:
		w.log.Info(fmt.Sprintf("Fetching ended, fetched %d results", len(results)))
	}

}

func (w *Worker) startFetching(ctx context.Context, entities []string) []<-chan processor.JobResult {
	resultChannels := make([]<-chan processor.JobResult, len(entities))
	for i, entity := range entities {
		ch := make(chan processor.JobResult, 1)

		fetchFn := processor.GetProcessEntityFn(entity)
		go func(fetchTweetsForEntity processor.ProcessEntityFn) {
			defer close(ch)
			select {
			case ch <- fetchTweetsForEntity(ctx):
			case <-ctx.Done():
				return
			}
		}(fetchFn)
		resultChannels[i] = ch
	}

	return resultChannels
}

func (w *Worker) collectResults(results *processor.JobResults, resultsChannels []<-chan processor.JobResult) chan struct{} {
	fetchingEnded := make(chan struct{})
	go func() {
		defer close(fetchingEnded)
		for _, ch := range resultsChannels {
			for result := range ch {
				if result.Error != nil {
					w.log.Error(fmt.Sprintf("Received error job result: %v", result.Error))
				}
				w.log.Info(fmt.Sprintf("Received result: %s", result.EntityId))
				*results = append(*results, result)
			}
		}
	}()

	return fetchingEnded
}
