package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Logger   LoggerConfig
	Postgres PostgresConfig
	Token    JWTConfig
	CORS     CORSConfig
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

	postgresConfig, err := newPostgresConfig()
	if err != nil {
		return nil, err
	}

	jwtCinfig, err := newJWTConfig()
	if err != nil {
		return nil, err
	}

	corsConfig, err := newCORSConfig()
	if err != nil {
		return nil, err
	}

	config := Config{
		Server:   serverConfig,
		Logger:   loggerConfig,
		Postgres: postgresConfig,
		Token:    jwtCinfig,
		CORS:     corsConfig,
	}

	return &config, nil
}
