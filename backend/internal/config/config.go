package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                               string
	AppPort                              string
	DatabaseURL                          string
	CircleAgentWalletExecutionEnabled    bool
	CircleAgentOnboardingOTPStartEnabled bool
	CircleAgentWalletFaucetEnabled       bool
	CircleCLIPath                        string
	CircleAgentWalletChain               string
	CircleAgentWalletTimeoutSeconds      int
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:                               os.Getenv("APP_ENV"),
		AppPort:                              os.Getenv("APP_PORT"),
		DatabaseURL:                          os.Getenv("DATABASE_URL"),
		CircleAgentWalletExecutionEnabled:    parseBool(os.Getenv("CIRCLE_AGENT_WALLET_EXECUTION_ENABLED")),
		CircleAgentOnboardingOTPStartEnabled: parseBool(os.Getenv("CIRCLE_AGENT_ONBOARDING_OTP_START_ENABLED")),
		CircleAgentWalletFaucetEnabled:       parseBool(os.Getenv("CIRCLE_AGENT_WALLET_FAUCET_ENABLED")),
		CircleCLIPath:                        os.Getenv("CIRCLE_CLI_PATH"),
		CircleAgentWalletChain:               os.Getenv("CIRCLE_AGENT_WALLET_CHAIN"),
		CircleAgentWalletTimeoutSeconds:      parseInt(os.Getenv("CIRCLE_AGENT_WALLET_TIMEOUT_SECONDS"), 120),
	}

	if cfg.AppEnv == "" {
		cfg.AppEnv = "development"
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "4000"
	}
	if cfg.CircleCLIPath == "" {
		cfg.CircleCLIPath = "circle"
	}
	if cfg.CircleAgentWalletChain == "" {
		cfg.CircleAgentWalletChain = "ARC-TESTNET"
	}

	return cfg
}

func (c Config) ValidateDatabaseURL() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	return nil
}

func parseBool(value string) bool {
	parsed, err := strconv.ParseBool(value)
	return err == nil && parsed
}

func parseInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
