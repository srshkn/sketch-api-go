package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"sketch-api-go/internal/config"
)

func NewPool(ctx context.Context, cfg config.Postgres) (*pgxpool.Pool, error) {
	url := cfg.URL()
	if url == "" {
		return nil, errors.New("database URL is empty")
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}

	config.MaxConns = cfg.MaxOpenConns()
	config.MinConns = cfg.MinOpenConns()
	config.MaxConnLifetime = cfg.ConnMaxLifetime()
	config.MaxConnIdleTime = cfg.MaxConnIdleTime()

	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, config)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return pool, nil
}
