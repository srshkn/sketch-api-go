package config

import (
	"fmt"
	"os"
)

const (
	loggerFormatEnv string = "LOGGER_FORMAT"
)

type LoggerConfig struct {
	Format string
}

func (l *LoggerConfig) validateLogger() error {

	if l.Format != "dev" && l.Format != "local" && l.Format != "prod" {
		return fmt.Errorf("environment variable %q is required", loggerFormatEnv)
	}

	return nil
}

func newLoggerConfig() (LoggerConfig, error) {
	logger := LoggerConfig{
		Format: os.Getenv(loggerFormatEnv),
	}

	if err := logger.validateLogger(); err != nil {
		return logger, err
	}

	return logger, nil
}
