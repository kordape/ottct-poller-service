package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/kordape/ottct-poller-service/config"
	"github.com/kordape/ottct-poller-service/internal/event"
	"github.com/kordape/ottct-poller-service/internal/processor"
	"github.com/kordape/ottct-poller-service/internal/worker"
	"github.com/kordape/ottct-poller-service/pkg/logger"
	"github.com/kordape/ottct-poller-service/pkg/predictor"
	"github.com/kordape/ottct-poller-service/pkg/sqs"
	"github.com/kordape/ottct-poller-service/pkg/twitter"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log := logger.New(cfg.Log.Level)

	awsConfig, err := initAWSConfig(cfg.FakeNewsQueue.SQSRegion, cfg.FakeNewsQueue.SQSAWSEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	awsSQSClient := awssqs.NewFromConfig(awsConfig)
	sqsClient := sqs.NewClient(awsSQSClient, cfg.FakeNewsQueue.SQSQueueURL)

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
		event.SendFakeNewsEventFnBuilder(sqsClient, log),
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

func initAWSConfig(region, endpoint string) (aws.Config, error) {
	if len(endpoint) > 0 {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			if service == awssqs.ServiceID {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}
			// Returning EndpointNotFoundError will allow the service to fallback
			// to it's default resolution.
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		return awsconfig.LoadDefaultConfig(
			context.Background(),
			awsconfig.WithEndpointResolverWithOptions(customResolver),
			awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("id", "fake-secret", "fake-token")),
			awsconfig.WithRegion(region),
		)
	}

	return awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(region))
}
