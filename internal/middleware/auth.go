package middleware

import (
	"context"
	"errors"
	"net/http"
	v1Generated "sketch-api-go/internal/generated/v1"
	"sketch-api-go/internal/token"
	"strings"
)

const (
	userIDKey string = "user_id"
)

func extractBearerToken(rowToken string) (string, error) {

	if rowToken == "" {
		return "", errors.New("missing authorization header")
	}

	const prefix = "Bearer "

	if !strings.HasPrefix(rowToken, prefix) {
		return "", errors.New("invalid authorization scheme")
	}

	token := strings.TrimSpace(strings.TrimPrefix(rowToken, prefix))
	if token == "" {
		return "", errors.New("missing bearer token")
	}

	return token, nil
}

func AuthMiddleware(tokenManager *token.Manager) v1Generated.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scopes := r.Context().Value(v1Generated.BearerAuthScopes)

			// Если operation не требует BearerAuth —
			// пропускаем запрос.
			if scopes == nil {
				next.ServeHTTP(w, r)
				return
			}

			// Здесь уже проверяем JWT.
			tokenString, err := extractBearerToken(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := tokenManager.ValidateAccessToken(tokenString)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				userIDKey,
				claims.UserID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
