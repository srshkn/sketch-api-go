package main

import (
	"log"
	"net"
	"net/http"

	"sketch-api-go/internal/config"
	"sketch-api-go/internal/generated"
	handler "sketch-api-go/internal/handlers"
	"sketch-api-go/internal/swagger"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	apiHandler := handler.New()

	generated.HandlerFromMux(apiHandler, mux)

	swagger.Register(mux)

	addr := net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)

	log.Printf("API: http://%s\n", addr)
	log.Printf("Swagger UI: http://%s/docs/\n", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
