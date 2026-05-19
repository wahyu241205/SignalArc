package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	AppPort     string
	DatabaseURL string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:      os.Getenv("APP_ENV"),
		AppPort:     os.Getenv("APP_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.AppEnv == "" {
		cfg.AppEnv = "development"
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "4000"
	}

	return cfg
}

func (c Config) ValidateDatabaseURL() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	return nil
}
