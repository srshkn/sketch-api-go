package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	serverHostEnv            string = "SERVER_HOST"
	serverPortEnv            string = "SERVER_PORT"
	shutdownTimeoutSecondEnv string = "SHOTDOWN_TIMEOUT_SECOND"
)

type Server interface {
	Host() string
	Port() string
	ShutdownTimeout() time.Duration
}

type configServer struct {
	host                  string
	port                  string
	shutdownTimeoutSecond time.Duration
}

func (s *configServer) Host() string {
	return s.host
}

func (s *configServer) Port() string {
	return s.port
}

func (s *configServer) ShutdownTimeout() time.Duration {
	return s.shutdownTimeoutSecond
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

	if s.shutdownTimeoutSecond <= 0 {
		return fmt.Errorf(
			"environment variable %q must be greater than 0",
			shutdownTimeoutSecondEnv,
		)
	}

	return nil
}

func (s *configServer) newShutdownTimeout() error {
	timeString := os.Getenv(shutdownTimeoutSecondEnv)
	if timeString == "" {
		return fmt.Errorf(
			"environment variable %q is required",
			shutdownTimeoutSecondEnv,
		)
	}

	timeout, err := strconv.Atoi(timeString)
	if err != nil {
		return err
	}

	s.shutdownTimeoutSecond = time.Duration(timeout)

	return nil
}

func newServerConfig() (*configServer, error) {

	server := configServer{
		host: os.Getenv(serverHostEnv),
		port: os.Getenv(serverPortEnv),
	}

	if err := server.newShutdownTimeout(); err != nil {
		return &server, err
	}

	if err := server.validateServer(); err != nil {
		return &server, err
	}

	return &server, nil
}
