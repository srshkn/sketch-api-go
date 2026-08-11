package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sketch-api-go/internal/app"
	"sketch-api-go/internal/config"
	"sketch-api-go/internal/db"
	"sketch-api-go/internal/logging"
	"sketch-api-go/internal/postgres"
	"sketch-api-go/internal/token"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		slog.Error(
			"failed to load configuration",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger := logging.New(cfg.Logger)

	pool, err := postgres.NewPool(ctx, cfg.Postgres)
	if err != nil {
		logger.Error(
			"connect to PostgreSQL: %v",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer pool.Close()

	jwtManager := token.New(cfg.Token)

	queries := db.New(pool)

	serverApp := app.New(
		cfg.Server,
		logger,
		queries,
		jwtManager,
		cfg.CORS,
	)

	err = serverApp.Run(ctx)
	if err != nil {
		logger.Error(
			"server stopped unexpectedly",
			slog.Any("error", err),
		)
	}
}
