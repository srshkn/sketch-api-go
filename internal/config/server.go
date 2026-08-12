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

type Server interface {
	Host() string
	Port() string
}

type configServer struct {
	host string
	port string
}

func (s *configServer) Host() string {
	return s.host
}

func (s *configServer) Port() string {
	return s.port
}

func (s *configServer) validateServer() error {
	switch {
	case s.host == "":
		return fmt.Errorf("environment variable %q is required", serverHostEnv)
	case s.port == "":
		return fmt.Errorf("environment variable %q is required", serverPortEnv)
	}

	port, err := strconv.Atoi(s.port)
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

func newServerConfig() (*configServer, error) {
	server := configServer{
		host: os.Getenv(serverHostEnv),
		port: os.Getenv(serverPortEnv),
	}

	if err := server.validateServer(); err != nil {
		return &server, err
	}

	return &server, nil
}
