package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App    `yaml:"app"`
		Log    `yaml:"logger"`
		Worker `yaml:"worker"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
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
