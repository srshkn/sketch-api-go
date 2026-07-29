package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	serverHostEnv string = "SERVER_HOST"
	serverPortEnv string = "SERVER_PORT"
)

type ServerConfig struct {
	Host string
	Port string
}

func (s *ServerConfig) validateServer() error {
	switch {
	case s.Host == "":
		return fmt.Errorf("environment variable %q is required", serverHostEnv)
	case s.Port == "":
		return fmt.Errorf("environment variable %q is required", serverPortEnv)
	}

	port, err := strconv.Atoi(s.Port)
	if err != nil {
		return fmt.Errorf(
			"environment variable %q must be a number: %w",
			serverPortEnv,
			err,
		)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf(
			"environment variable %q must be between 1 and 65535",
			serverPortEnv,
		)
	}

	return nil
}

func newServerConfig() (ServerConfig, error) {
	server := ServerConfig{
		Host: os.Getenv(serverHostEnv),
		Port: os.Getenv(serverPortEnv),
	}

	if err := server.validateServer(); err != nil {
		return server, err
	}

	return server, nil
}
