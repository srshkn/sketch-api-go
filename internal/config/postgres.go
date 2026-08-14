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

	maxOpenConnsEnv    string = "MAX_OPEN_CONNS"
	minOpenConnsEnv    string = "MIN_OPEN_CONNS"
	connMaxLifetimeEnv string = "CONN_MAX_LIFETIME"
	maxConnIdleTimeEnv string = "MAX_CONN_IDLE_TIME"
)

type Postgres interface {
	URL() string
	MaxOpenConns() int32
	MinOpenConns() int32
	ConnMaxLifetime() time.Duration
	MaxConnIdleTime() time.Duration
}

type configPostgres struct {
	host            string
	port            string
	user            string
	password        string
	dataBase        string
	url             string
	maxOpenConns    int32
	minOpenConns    int32
	connMaxLifetime time.Duration
	maxConnIdleTime time.Duration
}

func (p *configPostgres) URL() string {
	return p.url
}

func (p *configPostgres) MaxOpenConns() int32 {
	return p.maxOpenConns
}

func (p *configPostgres) MinOpenConns() int32 {
	return p.minOpenConns
}

func (p *configPostgres) ConnMaxLifetime() time.Duration {
	return p.connMaxLifetime
}

func (p *configPostgres) MaxConnIdleTime() time.Duration {
	return p.maxConnIdleTime
}

func (p *configPostgres) validatePostgres() error {
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

func (p *configPostgres) createUrl() {
	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.user, p.password),
		Host:   fmt.Sprintf("%s:%s", p.host, p.port),
		Path:   p.dataBase,
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	p.url = databaseURL.String()
}

func newPostgresConfig() (*configPostgres, error) {
	var postgres configPostgres

	maxOpenConnsStr := os.Getenv(maxOpenConnsEnv)
	if maxOpenConnsStr == "" {
		return &postgres, fmt.Errorf("")
	}

	maxOpenConns, err := strconv.Atoi(maxOpenConnsStr)
	if err != nil {
		return &postgres, fmt.Errorf(
			"environment variable %q is required",
			maxOpenConnsEnv,
		)
	}

	minOpenConnsStr := os.Getenv(minOpenConnsEnv)
	if minOpenConnsStr == "" {
		return &postgres, fmt.Errorf("")
	}

	minOpenConns, err := strconv.Atoi(minOpenConnsStr)
	if err != nil {
		return &postgres, fmt.Errorf(
			"environment variable %q is required",
			minOpenConnsEnv,
		)
	}

	connMaxLifetimeStr := os.Getenv(connMaxLifetimeEnv)
	if connMaxLifetimeStr == "" {
		return &postgres, fmt.Errorf("")
	}

	connMaxLifetime, err := strconv.Atoi(connMaxLifetimeStr)
	if err != nil {
		return &postgres, fmt.Errorf(
			"environment variable %q is required",
			connMaxLifetimeEnv,
		)
	}

	maxConnIdleTimeStr := os.Getenv(maxConnIdleTimeEnv)
	if maxConnIdleTimeStr == "" {
		return &postgres, fmt.Errorf("")
	}

	maxConnIdleTime, err := strconv.Atoi(maxConnIdleTimeStr)
	if err != nil {
		return &postgres, fmt.Errorf(
			"environment variable %q is required",
			maxConnIdleTimeEnv,
		)
	}

	postgres = configPostgres{
		host:            os.Getenv(postgresHostEnv),
		port:            os.Getenv(postgresPortEnv),
		user:            os.Getenv(postgresUserEnv),
		password:        os.Getenv(postgresPasswordEnv),
		dataBase:        os.Getenv(postgresDataBaseEnv),
		maxOpenConns:    int32(maxOpenConns),
		minOpenConns:    int32(minOpenConns),
		connMaxLifetime: time.Duration(connMaxLifetime) * time.Hour,
		maxConnIdleTime: time.Duration(maxConnIdleTime) * time.Minute,
	}

	if err := postgres.validatePostgres(); err != nil {
		return &postgres, err
	}

	postgres.createUrl()

	return &postgres, nil
}
