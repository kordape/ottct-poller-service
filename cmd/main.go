package main

import (
	"log"
	"time"

	"github.com/kordape/ottct-poller-service/config"
	"github.com/kordape/ottct-poller-service/internal/processor"
	"github.com/kordape/ottct-poller-service/internal/worker"
	"github.com/kordape/ottct-poller-service/pkg/logger"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log := logger.New(cfg.Log.Level)

	_, err = worker.NewWorker(
		log,
		processor.GetProcessEntityFn(),
		worker.WithInterval(time.Hour*time.Duration(cfg.IntervalHours)),
	)

	if err != nil {
		log.Fatal(err)
	}
}
