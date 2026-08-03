package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sketch-api-go/internal/generated"
	"sketch-api-go/internal/service"
	"strings"
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
	var request generated.RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			generated.INVALIDREQUEST,
			"invalid request body",
		)
		return
	}

	request.Username = strings.TrimSpace(request.Username)

	if request.Username == "" {
		writeError(
			w,
			http.StatusBadRequest,
			generated.INVALIDREQUEST,
			"name must not be empty",
		)
		return
	}

	if request.Password == "" {
		writeError(
			w,
			http.StatusBadRequest,
			generated.INVALIDREQUEST,
			"password must not be empty",
		)
		return
	}

	req, err := u.service.Registration(r.Context(), generated.RegisterUserRequest{
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
			generated.INVALIDREQUEST,
			"password must not be empty",
		)
		return
	}

	response := generated.UserResponse{
		Id:       int64(req.ID),
		Username: req.Name,
	}

	writeJSON(w, http.StatusCreated, response)
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
	code generated.ErrorResponseErrorCode,
	message string,
) {
	response := generated.ErrorResponse{}

	response.Error.Code = code
	response.Error.Message = message

	writeJSON(w, status, response)
}
