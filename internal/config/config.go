package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	serverHostEnv string = "SERVER_HOST"
	serverPortEnv string = "SERVER_PORT"
)

func getEnvError(name string) error {
	return fmt.Errorf("environment variable %q is required", name)
}

type ServerConfig struct {
	Host string
	Port string
}

func (s *ServerConfig) validateServer() error {
	switch {
	case s.Host == "":
		return getEnvError(serverHostEnv)
	case s.Port == "":
		return getEnvError(serverPortEnv)
	}

	return nil
}

type Config struct {
	Server ServerConfig
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	server := ServerConfig{
		Host: os.Getenv(serverHostEnv),
		Port: os.Getenv(serverPortEnv),
	}

	if err := server.validateServer(); err != nil {
		return nil, err
	}

	config := Config{
		Server: server,
	}

	return &config, nil
}
