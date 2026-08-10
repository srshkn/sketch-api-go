package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"sketch-api-go/internal/service"

	v1Generated "sketch-api-go/internal/generated/v1"
)

const (
	userIDKey string = "user_id"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (u *UserHandler) RegisterUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request v1Generated.RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			v1Generated.INVALIDREQUEST,
			"invalid request body",
		)
		return
	}

	request.Username = strings.TrimSpace(request.Username)

	if request.Username == "" {
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

	req, err := u.service.Registration(r.Context(), v1Generated.RegisterUserRequest{
		Confirmation: request.Confirmation,
		Email:        request.Email,
		Password:     request.Password,
		Username:     request.Username,
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

	writeJSON(w, http.StatusCreated, response)
}

func (u *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := u.service.GetUser(r.Context(), userID)
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
		Id:       user.ID,
		Username: user.Username,
	}

	writeJSON(w, http.StatusCreated, response)
}
