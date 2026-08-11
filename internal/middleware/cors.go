package middleware

import (
	"net/http"
	"sketch-api-go/internal/config"
	v1Generated "sketch-api-go/internal/generated/v1"
	"slices"
)

func CORS(cfg config.CORSConfig) v1Generated.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if slices.Contains(cfg.AllowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", cfg.AllowCredentials)
			}

			w.Header().Set(
				"Access-Control-Allow-Methods",
				cfg.AllowedMethods,
			)

			w.Header().Set(
				"Access-Control-Allow-Headers",
				cfg.AllowedHeaders,
			)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
