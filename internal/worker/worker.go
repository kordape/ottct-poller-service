package worker

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kordape/ottct-poller-service/internal/database"
	"github.com/kordape/ottct-poller-service/internal/event"
	"github.com/kordape/ottct-poller-service/internal/processor"
	"github.com/kordape/ottct-poller-service/pkg/logger"
)

const (
	defaultTickInterval         = 10 * time.Second
	defaultProcessorTimeoutInMs = 2 * int64(time.Millisecond)
	taskPoolSize                = 2
)

type Worker struct {
	tickInterval time.Duration
	log          logger.Interface

	running     int32
	stopChannel chan bool

	processorTimeoutInMs int64
	processor            processor.ProcessFn
	fakeNewsEventSender  event.SendFakeNewsEventFn
	entityStorage        database.EntityStorage
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

func NewWorker(log logger.Interface, processor processor.ProcessFn, fakeNewsEventSender event.SendFakeNewsEventFn, entityStorage database.EntityStorage, opts ...Option) (*Worker, error) {
	stopChan := make(chan bool)

	w := &Worker{
		tickInterval:         defaultTickInterval,
		log:                  log,
		stopChannel:          stopChan,
		processorTimeoutInMs: defaultProcessorTimeoutInMs,
		processor:            processor,
		fakeNewsEventSender:  fakeNewsEventSender,
		entityStorage:        entityStorage,
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

	if w.processor == nil {
		return errors.New("entity processor is nil")
	}

	if w.fakeNewsEventSender == nil {
		return errors.New("fake news event sender is nil")
	}

	if w.entityStorage == nil {
		return errors.New("entity storage is nil")
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
				go func() {
					// create processing task
					w.log.Info("Worker tick")
					results, err := w.process()
					if err != nil {
						w.log.Error(fmt.Sprintf("Processor finished with error: %v", err))
					}
					w.log.Info(fmt.Sprintf("Worker tick done, got %d results", len(results)))
					w.postProcess(results)
				}()

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
	ctx := context.Background()
	endTime := time.Now()
	startTime := endTime.Add(-w.tickInterval)

	w.log.Info(fmt.Sprintf("Processing for interval: %s - %s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339)))

	entities, err := w.entityStorage.GetEntities(ctx)
	if err != nil {
		return processor.JobResults{}, fmt.Errorf("failed to get entities: %w", err)
	}

	requests := make([]processor.JobRequest, len(entities))
	for i, e := range entities {
		requests[i] = processor.JobRequest{
			EntityID:  e.TwitterId,
			StartTime: startTime,
			EndTime:   endTime,
		}
	}

	results := w.pooledTasks(ctx, requests)
	return results, nil

}

func (w *Worker) pooledTasks(ctx context.Context, requests []processor.JobRequest) processor.JobResults {
	numJobs := len(requests)
	jobs := make(chan processor.JobRequest, numJobs)
	results := make(chan processor.JobResult, numJobs)
	defer close(results)

	for t := 0; t < taskPoolSize; t++ {
		go w.task(ctx, t, jobs, results)
	}

	for _, r := range requests {
		jobs <- r
	}
	close(jobs)

	//collectResults
	jobResults := make(processor.JobResults, numJobs)
	for i := 0; i < numJobs; i++ {
		jobResults[i] = <-results
	}

	return jobResults

}

func (w *Worker) task(ctx context.Context, id int, jobs <-chan processor.JobRequest, results chan<- processor.JobResult) {
	for job := range jobs {
		request := job
		select {
		case <-time.After(time.Millisecond * 2):
			w.log.Info("PROCESSOR TIMEOUT")
			results <- processor.JobResult{
				EntityID: job.EntityID,
				Error:    context.Canceled,
			}
		case results <- w.processor(ctx, request):
		}
	}
}

func (w *Worker) postProcess(results processor.JobResults) error {
	events := []event.FakeNews{}

	for _, result := range results {
		if result.Error != nil {
			w.log.Debug(fmt.Sprintf("Skipping error result: %s", result.Error))
			continue
		}

		for _, fakeNewsTweet := range result.FakeNewsTweets {
			events = append(events, event.FakeNews{
				EntityId:  result.EntityID,
				Timestamp: fakeNewsTweet.Timestamp,
				Content:   fakeNewsTweet.Content,
			})
		}
	}

	err := w.fakeNewsEventSender(context.Background(), events)

	if err != nil {
		return errors.New("failed to send events")
	}

	return nil
}
