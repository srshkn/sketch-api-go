package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sketch-api-go/internal/app"
	"sketch-api-go/internal/config"
	"sketch-api-go/internal/logging"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	slog.Info("read config...")

	cfg, err := config.New()
	if err != nil {
		slog.Error(
			"failed to load configuration",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	serverApp := app.New(cfg)

	logger := logging.New(cfg.Logger.Format)

	logger.Info("starting server",
		slog.String("address", serverApp.Server.Addr),
	)

	logger.Info("server endpoints",
		slog.String("api", "http://"+serverApp.Server.Addr),
		slog.String("swagger", "http://"+serverApp.Server.Addr+"/docs/"),
	)

	err = serverApp.Run(ctx)
	if err != nil {
		logger.Error(
			"server stopped unexpectedly",
			slog.Any("error", err),
		)
	}
}
