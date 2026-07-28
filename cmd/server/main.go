package main

import (
	"log"
	"net/http"

	"sketch-api-go/internal/generated"
	handler "sketch-api-go/internal/handlers"
	"sketch-api-go/internal/swagger"
)

func main() {
	mux := http.NewServeMux()

	apiHandler := handler.New()

	generated.HandlerFromMux(apiHandler, mux)

	swagger.Register(mux)

	log.Println("API: http://localhost:8080")
	log.Println("Swagger UI: http://localhost:8080/docs/")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
