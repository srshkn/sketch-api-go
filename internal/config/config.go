package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

/*
type ConfigApp interface {
	Server() *configServer
	Logger() *configLogger
	Postgres() *configPostgres
	JWT() *configJWT
	CORS() *configCORS
}
*/

type config struct {
	server   configServer
	logger   configLogger
	postgres configPostgres
	jwt      configJWT
	cors     configCORS
	cookie   configCookie
}

func (c *config) Server() *configServer {
	return &c.server
}

func (c *config) Logger() *configLogger {
	return &c.logger
}

func (c *config) Postgres() *configPostgres {
	return &c.postgres
}

func (c *config) JWT() *configJWT {
	return &c.jwt
}

func (c *config) CORS() *configCORS {
	return &c.cors
}

func (c *config) Cookie() *configCookie {
	return &c.cookie
}

func New() (*config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	server, err := newServerConfig()
	if err != nil {
		return nil, err
	}

	logger, err := newLoggerConfig()
	if err != nil {
		return nil, err
	}

	postgres, err := newPostgresConfig()
	if err != nil {
		return nil, err
	}

	jwt, err := newJWTConfig()
	if err != nil {
		return nil, err
	}

	cors, err := newCORSConfig()
	if err != nil {
		return nil, err
	}

	cookie, err := newCookieConfig()
	if err != nil {
		return nil, err
	}

	config := config{
		server:   *server,
		logger:   *logger,
		postgres: *postgres,
		jwt:      *jwt,
		cors:     *cors,
		cookie:   *cookie,
	}

	return &config, nil
}
