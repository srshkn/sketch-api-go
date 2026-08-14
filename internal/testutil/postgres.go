package testutil

import (
	"context"
	"fmt"
	"sketch-api-go/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTestPostgres(ctx context.Context, cfg config.Postgres, url string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf(
			"parse database URL: %w",
			err,
		)
	}

	poolConfig.MaxConns = cfg.MaxOpenConns()
	poolConfig.MinConns = cfg.MinOpenConns()
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime()
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime()

	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(
		connectCtx,
		poolConfig,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create PostgreSQL pool: %w",
			err,
		)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()

		return nil, fmt.Errorf(
			"ping PostgreSQL: %w",
			err,
		)
	}

	return pool, nil
}
