package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"sketch-api-go/internal/cookie"
	"sketch-api-go/internal/service"

	v1Generated "sketch-api-go/internal/generated/v1"
)

type AuthHandler struct {
	service *service.AuthService
	cookie  cookie.Auth
}

func NewAuthHandler(service *service.AuthService, cookie cookie.Auth) *AuthHandler {
	return &AuthHandler{
		service: service,
		cookie:  cookie,
	}
}

func (a *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request v1Generated.LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			v1Generated.INVALIDREQUEST,
			"invalid request body",
		)
		return
	}

	if request.Email == "" {
		writeError(
			w,
			http.StatusBadRequest,
			v1Generated.INVALIDREQUEST,
			"name must not be empty",
		)
		return
	}

	if request.Password == "" {
		writeError(
			w,
			http.StatusBadRequest,
			v1Generated.INVALIDREQUEST,
			"password must not be empty",
		)
		return
	}

	tokenPair, err := a.service.Login(r.Context(), request)
	if err != nil {
		slog.Error(
			"user registration failed",
			slog.Any("error", err),
		)

		writeError(
			w,
			http.StatusBadRequest,
			v1Generated.INVALIDREQUEST,
			"password must not be empty",
		)
		return
	}

	response := v1Generated.TokensResponse{
		AccessToken:  tokenPair.GetAccessToken(),
		RefreshToken: tokenPair.GetRefreshToken(),
	}

	a.cookie.SetRefreshToken(w, tokenPair.GetRefreshToken(), tokenPair.ExpiresAt)

	writeJSON(w, http.StatusOK, response)
}

func (a *AuthHandler) UpdateRefreshToken(w http.ResponseWriter, r *http.Request) {
	var request v1Generated.RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			v1Generated.INVALIDREQUEST,
			"invalid request body",
		)
		return
	}

	if request.RefreshToken == "" {
		writeError(
			w,
			http.StatusBadRequest,
			v1Generated.INVALIDREQUEST,
			"name must not be empty",
		)
		return
	}

	tokenPair, err := a.service.Refresh(r.Context(), request)
	if err != nil {
		slog.Error(
			"user registration failed",
			slog.Any("error", err),
		)

		writeError(
			w,
			http.StatusBadRequest,
			v1Generated.INVALIDREQUEST,
			"...",
		)
		return
	}

	response := v1Generated.TokensResponse{
		AccessToken:  tokenPair.GetAccessToken(),
		RefreshToken: tokenPair.GetRefreshToken(),
	}

	a.cookie.SetRefreshToken(w, tokenPair.GetRefreshToken(), tokenPair.ExpiresAt)

	writeJSON(w, http.StatusOK, response)
}
