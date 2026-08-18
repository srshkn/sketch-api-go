package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"sketch-api-go/internal/config"
	"sketch-api-go/internal/cookie"
	"sketch-api-go/internal/db"
	"sketch-api-go/internal/service"
	"sketch-api-go/internal/swagger"

	v1Generated "sketch-api-go/internal/generated/v1"
	v1Handler "sketch-api-go/internal/handler/v1"
	"sketch-api-go/internal/middleware"
	"sketch-api-go/internal/token"
)

type serverApp struct {
	server          *http.Server
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

func New(
	cfg config.Server,
	logger *slog.Logger,
	db db.Querier,
	jwtManager token.JWTManager,
	cors config.CORS,
	cookieManager cookie.Auth,
) *serverApp {
	mux := http.NewServeMux()

	userService := service.NewUserService(db)
	authService := service.NewAuthService(db, jwtManager)

	v1MetaHandler := v1Handler.NewMetaHandler()
	v1UserHandler := v1Handler.NewUserHandler(userService)
	v1AuthHandler := v1Handler.NewAuthHandler(authService, cookieManager)

	v1APIHandler := v1Handler.New(
		v1MetaHandler,
		v1UserHandler,
		v1AuthHandler,
	)

	v1Generated.HandlerWithOptions(
		v1APIHandler,
		v1Generated.StdHTTPServerOptions{
			BaseRouter: mux,
			BaseURL:    "/api/v1",
			Middlewares: []v1Generated.MiddlewareFunc{
				middleware.Logging(logger),
				middleware.AuthMiddleware(jwtManager),
			},
		},
	)

	swagger.Register(mux)

	server := &http.Server{
		Handler: middleware.CORS(cors)(mux),
		Addr:    net.JoinHostPort(cfg.Host(), cfg.Port()),
	}

	return &serverApp{
		server:          server,
		logger:          logger,
		shutdownTimeout: cfg.ShutdownTimeout() * time.Second,
	}
}

func (s *serverApp) gracefulStop(serverErr <-chan error) error {
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		s.shutdownTimeout,
	)
	defer cancel()

	s.logger.Info(
		"shutting down HTTP server",
		slog.Duration("timeout", s.shutdownTimeout),
	)

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		s.logger.Error(
			"graceful shutdown failed",
			slog.Any("error", err),
		)

		if closeErr := s.server.Close(); closeErr != nil {
			s.logger.Error(
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

	s.logger.Info("HTTP server stopped")

	return nil
}

func (s *serverApp) Run(ctx context.Context) error {
	serverErr := make(chan error, 1)

	go func() {
		serverErr <- s.server.ListenAndServe()
	}()

	s.logger.Info("starting server",
		slog.String("address", s.server.Addr),
	)

	s.logger.Info("server endpoints",
		slog.String("api", "http://"+s.server.Addr),
		slog.String("swagger", "http://"+s.server.Addr+"/docs/"),
	)

	select {

	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err

	case <-ctx.Done():
		s.logger.Info("shutdown signal received")
		return s.gracefulStop(serverErr)

	}
}

func (s *serverApp) Handler() http.Handler {
	return s.server.Handler
}
