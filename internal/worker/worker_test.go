package worker

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/kordape/ottct-poller-service/internal/database"
	"github.com/kordape/ottct-poller-service/internal/event"
	"github.com/kordape/ottct-poller-service/internal/processor"
	"github.com/kordape/ottct-poller-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorker(t *testing.T) {

	t.Run("single tick all results processed", func(t *testing.T) {
		log := logger.New("DEBUG")

		processEntityFn := func(ctx context.Context, request processor.JobRequest) processor.JobResult {
			fakeNewsTweets := make([]processor.FakeNewsTweet, 10)
			for i := range fakeNewsTweets {
				fakeNewsTweets[i] = processor.FakeNewsTweet{
					Timestamp: time.Now(),
					Content:   fmt.Sprintf("Tweet%d", i),
				}
			}

			return processor.JobResult{
				EntityID:       request.EntityID,
				FakeNewsTweets: fakeNewsTweets,
			}
		}

		eventSenderFn := func(ctx context.Context, events []event.FakeNews) error {
			assert.Equal(t, 20, len(events))
			return nil
		}

		db := database.NewMockEntityStorage(t)
		db.On("GetEntities", mock.Anything).Return([]database.Entity{
			{
				ID:          "id1",
				TwitterId:   "foo",
				DisplayName: "foo",
			},
			{
				ID:          "id2",
				TwitterId:   "bar",
				DisplayName: "bar",
			},
		}, nil)

		w, err := NewWorker(log, processEntityFn, eventSenderFn, db, WithInterval(5*time.Second))
		assert.NoError(t, err)

		err = w.Run()
		assert.NoError(t, err)

		time.Sleep(8 * time.Second)
		w.Stop()
	})

	t.Run("single tick half results failed", func(t *testing.T) {
		log := logger.New("DEBUG")

		processEntityFn := func(ctx context.Context, request processor.JobRequest) processor.JobResult {
			if request.EntityID == "foo" {
				return processor.JobResult{
					EntityID: request.EntityID,
					Error:    errors.New("big error"),
				}
			}

			fakeNewsTweets := make([]processor.FakeNewsTweet, 10)
			for i := range fakeNewsTweets {
				fakeNewsTweets[i] = processor.FakeNewsTweet{
					Timestamp: time.Now(),
					Content:   fmt.Sprintf("Tweet%d", i),
				}
			}

			return processor.JobResult{
				EntityID:       request.EntityID,
				FakeNewsTweets: fakeNewsTweets,
			}
		}

		eventSenderFn := func(ctx context.Context, events []event.FakeNews) error {
			assert.Equal(t, 10, len(events))
			return nil
		}

		db := database.NewMockEntityStorage(t)
		db.On("GetEntities", mock.Anything).Return([]database.Entity{
			{
				ID:          "id1",
				TwitterId:   "foo",
				DisplayName: "foo",
			},
			{
				ID:          "id2",
				TwitterId:   "bar",
				DisplayName: "bar",
			},
		}, nil)

		w, err := NewWorker(log, processEntityFn, eventSenderFn, db, WithInterval(5*time.Second))
		assert.NoError(t, err)

		err = w.Run()
		assert.NoError(t, err)

		time.Sleep(8 * time.Second)
		w.Stop()
	})

	t.Run("single tick processing timeout", func(t *testing.T) {
		log := logger.New("DEBUG")

		processEntityFn := func(ctx context.Context, request processor.JobRequest) processor.JobResult {
			time.Sleep(20 * time.Second)

			return processor.JobResult{
				EntityID: request.EntityID,
			}
		}

		eventSenderFn := func(ctx context.Context, events []event.FakeNews) error {
			assert.Equal(t, 0, len(events))
			return nil
		}

		db := database.NewMockEntityStorage(t)
		db.On("GetEntities", mock.Anything).Return([]database.Entity{
			{
				ID:          "id1",
				TwitterId:   "foo",
				DisplayName: "foo",
			},
			{
				ID:          "id2",
				TwitterId:   "bar",
				DisplayName: "bar",
			},
		}, nil)

		w, err := NewWorker(log, processEntityFn, eventSenderFn, db, WithInterval(5*time.Second), WithProcessorTimeout(2000))
		assert.NoError(t, err)

		err = w.Run()
		assert.NoError(t, err)

		time.Sleep(8 * time.Second)
		w.Stop()
	})
}
