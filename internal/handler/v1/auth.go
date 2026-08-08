package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"sketch-api-go/internal/service"

	v1Generated "sketch-api-go/internal/generated/v1"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (a *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request v1Generated.LoginUserJSONRequestBody

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

	req, err := a.service.Login(r.Context(), v1Generated.LoginUserRequest{
		Email:    request.Email,
		Password: request.Password,
	})

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

	response := v1Generated.UserResponse{
		Id:       req.ID,
		Username: req.Username,
	}

	writeJSON(w, http.StatusOK, response)

}
