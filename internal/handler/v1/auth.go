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
	service service.Auth
	cookie  cookie.Auth
}

func NewAuthHandler(service service.Auth, cookie cookie.Auth) *AuthHandler {
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

	response, refreshToken, err := a.service.Login(r.Context(), request)
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

	a.cookie.SetRefreshToken(w, refreshToken.Token, refreshToken.ExpiresAt)
	writeJSON(w, http.StatusOK, response)
}

func (a *AuthHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("refresh_token")
	if err != nil {
		writeJSON(w, http.StatusNoContent, nil)
		return
	}

	err = a.service.Logout(r.Context(), token.Value)
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

	a.cookie.ClearRefreshToken(w)

	writeJSON(w, http.StatusNoContent, nil)
}

func (a *AuthHandler) UpdateRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(a.cookie.Name())
	if err != nil || cookie.Value == "" {
		writeError(
			w,
			http.StatusUnauthorized,
			v1Generated.INVALIDREQUEST,
			"refresh token is missing",
		)
		return
	}

	response, refreshToken, err := a.service.Refresh(r.Context(), cookie.Value)
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

	a.cookie.SetRefreshToken(w, refreshToken.Token, refreshToken.ExpiresAt)
	writeJSON(w, http.StatusOK, response)
}
