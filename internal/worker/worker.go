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

const (
	defaultTickInterval         = 1 * time.Second
	defaultProcessorTimeoutInMs = 2 * int64(time.Millisecond)
)

type Worker struct {
	tickInterval time.Duration
	log          logger.Interface

	running     int32
	stopChannel chan bool

	processorTimeoutInMs int64
	entityProcessor      processor.ProcessEntityFn
}

type Option func(w *Worker)

func WithInterval(interval time.Duration) Option {
	return func(w *Worker) {
		w.tickInterval = interval
	}
}

func WithProcessorTimeout(timeoutInMs int64) Option {
	return func(w *Worker) {
		w.processorTimeoutInMs = timeoutInMs
	}
}

func NewWorker(log logger.Interface, entityProcessor processor.ProcessEntityFn, opts ...Option) (*Worker, error) {
	stopChan := make(chan bool)

	w := &Worker{
		tickInterval:         defaultTickInterval,
		log:                  log,
		stopChannel:          stopChan,
		processorTimeoutInMs: defaultProcessorTimeoutInMs,
		entityProcessor:      entityProcessor,
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

	if w.entityProcessor == nil {
		return errors.New("entity processor is nil")
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
				results, err := w.process()
				if err != nil {
					w.log.Error(fmt.Sprintf("Processor finished with error: %v", err))
				}
				w.log.Info(fmt.Sprintf("Worker tick done, got %d results", len(results)))
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

func (w *Worker) process() (processor.JobResults, error) {
	endTime := time.Now()
	startTime := endTime.Add(-w.tickInterval)

	w.log.Info(fmt.Sprintf("Processing for interval: %s - %s", startTime.Format(time.UnixDate), endTime.Format(time.UnixDate)))

	// Create a context with a timeout
	ctxProcessor, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(w.processorTimeoutInMs))
	defer cancel()

	// TODO: replace this line with fetched entities from db
	entities := []string{"foo", "bar"}

	resultsChannels := w.startProcessing(ctxProcessor, entities)
	results := make(processor.JobResults, 0, len(entities))
	processingEnded := w.collectResults(&results, resultsChannels)

	select {
	case <-processingEnded:
		w.log.Info(fmt.Sprintf("Processing ended, processed %d results", len(results)))
		return results, nil
	case <-ctxProcessor.Done():
		// If context is cancelled (i.e. timeout reached)
		// return context canceled error
		return nil, context.Canceled
	}

}

func (w *Worker) startProcessing(ctx context.Context, entities []string) []<-chan processor.JobResult {
	resultChannels := make([]<-chan processor.JobResult, len(entities))
	for i, e := range entities {
		ch := make(chan processor.JobResult, 1)

		entity := e
		go func(processEntity processor.ProcessEntityFn) {
			defer close(ch)
			select {
			case ch <- processEntity(ctx, entity):
			case <-ctx.Done():
				return
			}
		}(w.entityProcessor)
		resultChannels[i] = ch
	}

	return resultChannels
}

func (w *Worker) collectResults(results *processor.JobResults, resultsChannels []<-chan processor.JobResult) chan struct{} {
	processingEnded := make(chan struct{})
	go func() {
		defer close(processingEnded)
		for _, ch := range resultsChannels {
			for result := range ch {
				if result.Error != nil {
					w.log.Error(fmt.Sprintf("Received error job result: %v", result.Error))
					continue
				}
				w.log.Info(fmt.Sprintf("Received result: %s", result.EntityId))
				*results = append(*results, result)
			}
		}
	}()

	return processingEnded
}
