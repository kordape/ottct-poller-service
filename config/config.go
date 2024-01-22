package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App           `yaml:"app"`
		Log           `yaml:"logger"`
		Worker        `yaml:"worker"`
		FakeNewsQueue `yaml:"fake_news_queue"`
		DB            `yaml:"db"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// Entities DB
	DB struct {
		URL string `env-required:"true" yaml:"username" env:"DB_URL"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	// Worker -.
	Worker struct {
		IntervalSeconds    int    `env-required:"true" yaml:"interval_seconds" env:"WORKER_INTERVAL_SECONDS"`
		TwitterBearerToken string `env-required:"true" yaml:"twitter_bearer_token" env:"TWITTER_BEARER_TOKEN"`
		PredictorBaseURL   string `env-required:"true" yaml:"predictor_base_url" env:"PREDICTOR_BASE_URL"`
	}

	// FakeNewsQueue holds configuration for `FakeNewsQueue` queue.
	FakeNewsQueue struct {
		SQSQueueURL    string `env-required:"true" yaml:"queue_url" env:"FAKE_NEWS_QUEUE_URL"`
		SQSAWSEndpoint string `yaml:"queue_endpoint" env:"FAKE_NEWS_QUEUE_ENDPOINT"`
		SQSRegion      string `env-required:"true" yaml:"queue_region" env:"FAKE_NEWS_QUEUE_REGION"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
