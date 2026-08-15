package v1

import (
	"encoding/json"
	"errors"
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
	service service.User
}

func NewUserHandler(service service.User) *UserHandler {
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
			"username must not be empty",
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

	if request.Password != request.Confirmation {
		writeError(
			w,
			http.StatusBadRequest,
			v1Generated.INVALIDREQUEST,
			"password and confirmation do not match",
		)
		return
	}

	response, err := u.service.Registration(r.Context(), v1Generated.RegisterUserRequest{
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
		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			writeError(
				w,
				http.StatusConflict,
				v1Generated.NOTFOUND,
				"a user with that name or email already exists",
			)
		default:
			slog.Error(
				"user registration failed",
				slog.Any("error", err),
			)

			writeError(
				w,
				http.StatusInternalServerError,
				v1Generated.INTERNALERROR,
				"internal server error",
			)
		}

		return
	}

	writeJSON(w, http.StatusCreated, response)
}

func (u *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := u.service.GetUser(r.Context(), userID)
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

	writeJSON(w, http.StatusCreated, response)
}
