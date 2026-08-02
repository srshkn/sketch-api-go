package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	postgresHostEnv     string = "POSTGRES_HOST"
	postgresPortEnv     string = "POSTGRES_PORT"
	postgresUserEnv     string = "POSTGRES_USER"
	postgresPasswordEnv string = "POSTGRES_PASSWORD"
	postgresDataBaseEnv string = "POSTGRES_DB"
)

type PostgresConfig struct {
	host            string
	port            string
	user            string
	password        string
	dataBase        string
	URL             string
	MaxOpenConns    int32
	MinOpenConns    int32
	ConnMaxLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func (p *PostgresConfig) validatePostgres() error {
	switch {
	case p.host == "":
		return fmt.Errorf("environment variable %q is required", postgresHostEnv)
	case p.port == "":
		return fmt.Errorf("environment variable %q is required", postgresPortEnv)
	case p.user == "":
		return fmt.Errorf("environment variable %q is required", postgresUserEnv)
	case p.password == "":
		return fmt.Errorf("environment variable %q is required", postgresPasswordEnv)
	case p.dataBase == "":
		return fmt.Errorf("environment variable %q is required", postgresDataBaseEnv)
	}

	port, err := strconv.Atoi(p.port)
	if err != nil {
		return fmt.Errorf(
			"environment variable %q must be a number: %w",
			postgresPortEnv,
			err,
		)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf(
			"environment variable %q must be between 1 and 65535",
			postgresPortEnv,
		)
	}

	if len(p.user) > 63 {
		return fmt.Errorf(
			"environment variable %q must not exceed 63 bytes",
			postgresUserEnv,
		)
	}

	if len(p.dataBase) > 63 {
		return fmt.Errorf(
			"environment variable %q must not exceed 63 bytes",
			postgresDataBaseEnv,
		)
	}

	return nil
}

func (p *PostgresConfig) createUrl() {
	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.user, p.password),
		Host:   fmt.Sprintf("%s:%s", p.host, p.port),
		Path:   p.dataBase,
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	p.URL = databaseURL.String()
}

func newPostgresConfig() (PostgresConfig, error) {
	postgres := PostgresConfig{
		host:            os.Getenv(postgresHostEnv),
		port:            os.Getenv(postgresPortEnv),
		user:            os.Getenv(postgresUserEnv),
		password:        os.Getenv(postgresPasswordEnv),
		dataBase:        os.Getenv(postgresDataBaseEnv),
		MaxOpenConns:    10,
		MinOpenConns:    2,
		ConnMaxLifetime: time.Hour,
		MaxConnIdleTime: 15 * time.Minute,
	}

	if err := postgres.validatePostgres(); err != nil {
		return postgres, err
	}

	postgres.createUrl()

	return postgres, nil
}
