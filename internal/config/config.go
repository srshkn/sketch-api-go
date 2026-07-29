package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Logger LoggerConfig
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	serverConfig, err := newServerConfig()
	if err != nil {
		return nil, err
	}

	loggerConfig, err := newLoggerConfig()
	if err != nil {
		return nil, err
	}

	config := Config{
		Server: serverConfig,
		Logger: loggerConfig,
	}

	return &config, nil
}
