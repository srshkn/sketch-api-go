package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	serverConfig, err := newServerConfig()
	if err != nil {
		return nil, err
	}

	config := Config{
		Server: serverConfig,
	}

	return &config, nil
}
