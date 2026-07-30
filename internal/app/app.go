package app

import (
	"net"
	"net/http"
	"sketch-api-go/internal/config"
	"sketch-api-go/internal/generated"
	"sketch-api-go/internal/handler"
	"sketch-api-go/internal/swagger"
)

type ServerApp struct {
	Server *http.Server
}

func New(cfg *config.Config) *ServerApp {

	app := new(ServerApp)

	mux := http.NewServeMux()

	apiHandler := handler.New()

	generated.HandlerFromMux(apiHandler, mux)

	swagger.Register(mux)

	addr := net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)

	app.Server = &http.Server{
		Handler: mux,
		Addr:    addr,
	}

	return app
}

func (s ServerApp) Run() error {
	err := s.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
