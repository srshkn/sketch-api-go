package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"sketch-api-go/internal/service"

	v1Generated "sketch-api-go/internal/generated/v1"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func writeError(
	w http.ResponseWriter,
	status int,
	code v1Generated.ErrorResponseErrorCode,
	message string,
) {
	response := v1Generated.ErrorResponse{}

	response.Error.Code = code
	response.Error.Message = message

	writeJSON(w, status, response)
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
		Username: req.Name,
	}

	writeJSON(w, http.StatusCreated, response)
}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
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

	req, err := u.service.Login(r.Context(), v1Generated.LoginUserRequest{
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
		Username: req.Name,
	}

	writeJSON(w, http.StatusOK, response)

}
