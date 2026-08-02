package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sketch-api-go/internal/config"
	"sketch-api-go/internal/generated"
	"sketch-api-go/internal/handler"
	"sketch-api-go/internal/swagger"
	"time"
)

type ServerApp struct {
	Server          *http.Server
	shutdownTimeout time.Duration
}

func New(cfg config.ServerConfig) *ServerApp {
	mux := http.NewServeMux()

	apiHandler := handler.New()
	generated.HandlerFromMux(apiHandler, mux)

	swagger.Register(mux)

	server := &http.Server{
		Handler: mux,
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
	}

	return &ServerApp{
		Server:          server,
		shutdownTimeout: 10 * time.Second,
	}
}

func (s *ServerApp) gracefulStop(serverErr <-chan error) error {
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)
	defer cancel()

	slog.Info(
		"shutting down HTTP server",
		slog.Duration("timeout", s.shutdownTimeout),
	)

	if err := s.Server.Shutdown(shutdownCtx); err != nil {
		slog.Error(
			"graceful shutdown failed",
			slog.Any("error", err),
		)

		if closeErr := s.Server.Close(); closeErr != nil {
			slog.Error(
				"force close HTTP server failed",
				slog.Any("error", closeErr),
			)
		}

		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	err := <-serverErr
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server stopped: %w", err)
	}

	slog.Info("HTTP server stopped")

	return nil
}

func (s *ServerApp) Run(ctx context.Context) error {
	serverErr := make(chan error, 1)

	go func() {
		serverErr <- s.Server.ListenAndServe()
	}()

	select {

	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err

	case <-ctx.Done():
		slog.Info("shutdown signal received")
		return s.gracefulStop(serverErr)

	}
}
