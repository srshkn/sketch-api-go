package handler

import (
	"encoding/json"
	"net/http"
	"sketch-api-go/internal/generated"
	"strings"
)

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
