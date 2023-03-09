package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kordape/ottct-poller-service/config"
	"github.com/kordape/ottct-poller-service/internal/event"
	"github.com/kordape/ottct-poller-service/internal/predictor"
	"github.com/kordape/ottct-poller-service/internal/processor"
	"github.com/kordape/ottct-poller-service/internal/twitter"
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

	w, err := worker.NewWorker(
		log,
		processor.GetProcessFn(
			log,
			twitter.New(
				&http.Client{
					Timeout: 10 * time.Second,
				},
				cfg.Worker.TwitterBearerToken,
			),
			predictor.New(
				&http.Client{
					Timeout: 10 * time.Second,
				},
				cfg.Worker.PredictorBaseURL,
			),
		),
		event.SendFakeNewsEventFnBuilder(),
		worker.WithInterval(time.Second*time.Duration(cfg.IntervalSeconds)),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = w.Run()

	if err != nil {
		log.Fatal(err)
	}

	// Wait for terminal signal.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	log.Info("Stopping worker")
	w.Stop()
}
