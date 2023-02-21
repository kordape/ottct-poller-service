package worker

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kordape/ottct-poller-service/internal/event"
	"github.com/kordape/ottct-poller-service/internal/processor"
	"github.com/kordape/ottct-poller-service/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {

	t.Run("single tick all results processed", func(t *testing.T) {
		log := logger.New("DEBUG")

		processEntityFn := func(ctx context.Context, entityId string) processor.JobResult {
			fakeNewsTweets := make([]processor.FakeNewsTweet, 10)
			for i := range fakeNewsTweets {
				fakeNewsTweets[i] = processor.FakeNewsTweet{
					Timestamp: time.Now().Unix(),
					Content:   fmt.Sprintf("Tweet%d", i),
				}
			}

			return processor.JobResult{
				EntityId:       entityId,
				FakeNewsTweets: fakeNewsTweets,
			}
		}

		eventSenderFn := func(events []event.FakeNews) error {
			assert.Equal(t, 20, len(events))
			return nil
		}

		w, err := NewWorker(log, processEntityFn, eventSenderFn, WithInterval(5*time.Second))
		assert.NoError(t, err)

		err = w.Run()
		assert.NoError(t, err)

		time.Sleep(8 * time.Second)
		w.Stop()
	})

	t.Run("single tick half results failed", func(t *testing.T) {
		log := logger.New("DEBUG")

		processEntityFn := func(ctx context.Context, entityId string) processor.JobResult {
			if entityId == "foo" {
				return processor.JobResult{
					EntityId: entityId,
					Error:    errors.New("big error"),
				}
			}

			fakeNewsTweets := make([]processor.FakeNewsTweet, 10)
			for i := range fakeNewsTweets {
				fakeNewsTweets[i] = processor.FakeNewsTweet{
					Timestamp: time.Now().Unix(),
					Content:   fmt.Sprintf("Tweet%d", i),
				}
			}

			return processor.JobResult{
				EntityId:       entityId,
				FakeNewsTweets: fakeNewsTweets,
			}
		}

		eventSenderFn := func(events []event.FakeNews) error {
			assert.Equal(t, 10, len(events))
			return nil
		}

		w, err := NewWorker(log, processEntityFn, eventSenderFn, WithInterval(5*time.Second))
		assert.NoError(t, err)

		err = w.Run()
		assert.NoError(t, err)

		time.Sleep(8 * time.Second)
		w.Stop()
	})

	t.Run("single tick processing timeout", func(t *testing.T) {
		log := logger.New("DEBUG")

		processEntityFn := func(ctx context.Context, entityId string) processor.JobResult {
			time.Sleep(20 * time.Second)

			return processor.JobResult{
				EntityId: entityId,
			}
		}

		eventSenderFn := func(events []event.FakeNews) error {
			assert.Equal(t, 0, len(events))
			return nil
		}

		w, err := NewWorker(log, processEntityFn, eventSenderFn, WithInterval(5*time.Second), WithProcessorTimeout(2000))
		assert.NoError(t, err)

		err = w.Run()
		assert.NoError(t, err)

		time.Sleep(8 * time.Second)
		w.Stop()
	})
}
