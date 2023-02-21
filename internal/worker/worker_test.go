package worker

import (
	"testing"
	"time"

	"github.com/kordape/ottct-poller-service/internal/processor"
	"github.com/kordape/ottct-poller-service/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {

	t.Run("worker tick", func(t *testing.T) {
		log := logger.New("DEBUG")
		w, err := NewWorker(log, processor.GetProcessEntityFn())
		assert.NoError(t, err)

		err = w.Run()
		assert.NoError(t, err)

		go func() {
			assert.True(t, w.Running())
			time.Sleep(5 * time.Second)
			w.Stop()
		}()

		time.Sleep(10 * time.Second)
		assert.False(t, w.Running())
	})

	// t.Run("worker - processor error", func(t *testing.T) {
	// 	log := logger.New("DEBUG")

	// 	processEntityFn := func(ctx context.Context, entityId string) processor.JobResult {
	// 		return processor.JobResult{
	// 			Error: errors.New("test error"),
	// 		}
	// 	}

	// 	w, err := NewWorker(log, processEntityFn)
	// 	assert.NoError(t, err)

	// 	err = w.Run()
	// 	assert.NoError(t, err)

	// 	go func() {
	// 		assert.True(t, w.Running())
	// 		time.Sleep(5 * time.Second)
	// 		w.Stop()
	// 	}()

	// 	time.Sleep(10 * time.Second)
	// 	assert.False(t, w.Running())
	// })

	// t.Run("worker - processor timeout", func(t *testing.T) {
	// 	log := logger.New("DEBUG")

	// 	processEntityFn := func(ctx context.Context, entityId string) processor.JobResult {
	// 		time.Sleep(100 * time.Millisecond)
	// 		return processor.JobResult{
	// 			EntityId: entityId,
	// 		}
	// 	}

	// 	w, err := NewWorker(log, processEntityFn, WithProcessorTimeout(80))
	// 	assert.NoError(t, err)

	// 	err = w.Run()
	// 	assert.NoError(t, err)

	// 	go func() {
	// 		assert.True(t, w.Running())
	// 		time.Sleep(5 * time.Second)
	// 		w.Stop()
	// 	}()

	// 	time.Sleep(10 * time.Second)
	// 	assert.False(t, w.Running())
	// })

}
