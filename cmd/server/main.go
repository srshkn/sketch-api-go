package main

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"

	"sketch-api-go/internal/config"
	"sketch-api-go/internal/generated"
	"sketch-api-go/internal/handler"
	"sketch-api-go/internal/logging"
	"sketch-api-go/internal/swagger"
)

func main() {
	slog.Info("read config...")

	cfg, err := config.New()
	if err != nil {
		slog.Error(
			"failed to load configuration",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger := logging.New(cfg.Logger.Format)

	mux := http.NewServeMux()

	apiHandler := handler.New()

	generated.HandlerFromMux(apiHandler, mux)

	swagger.Register(mux)

	addr := net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)

	logger.Info("starting server",
		slog.String("address", addr),
	)

	logger.Info("server endpoints",
		slog.String("api", "http://"+addr),
		slog.String("swagger", "http://"+addr+"/docs/"),
	)

	if err := http.ListenAndServe(addr, mux); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		logger.Error(
			"server stopped unexpectedly",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}
