package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"sketch-api-go/internal/generated"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

var _ generated.ServerInterface = (*Handler)(nil)

func (h *Handler) GetHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("OK"))
}

func (h *Handler) RegisterUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request generated.RegisterUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			generated.INVALIDREQUEST,
			"invalid request body",
		)
		return
	}

	request.Name = strings.TrimSpace(request.Name)

	if request.Name == "" {
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

	response := generated.UserResponse{
		Id:   1,
		Name: request.Name,
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
