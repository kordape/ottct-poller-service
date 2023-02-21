package worker

import (
	"testing"
	"time"

	"github.com/kordape/ottct-poller-service/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {

	t.Run("worker tick", func(t *testing.T) {
		log := logger.New("DEBUG")
		w, err := NewWorker(log)
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
}
