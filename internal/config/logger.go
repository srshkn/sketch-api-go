package config

import (
	"fmt"
	"os"
)

const (
	loggerFormatEnv string = "LOGGER_FORMAT"
)

type Logger interface {
	Format() string
}

type configLogger struct {
	format string
}

func (l *configLogger) Format() string {
	return l.format
}

func (l *configLogger) validateLogger() error {

	if l.format != "dev" && l.format != "local" && l.format != "prod" {
		return fmt.Errorf("environment variable %q is required", loggerFormatEnv)
	}

	return nil
}

func newLoggerConfig() (*configLogger, error) {
	logger := configLogger{
		format: os.Getenv(loggerFormatEnv),
	}

	if err := logger.validateLogger(); err != nil {
		return &logger, err
	}

	return &logger, nil
}
