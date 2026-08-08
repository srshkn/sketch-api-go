package v1

import (
	"encoding/json"
	"net/http"

	v1Generated "sketch-api-go/internal/generated/v1"
)

type Handler struct {
	meta *MetaHandler
	user *UserHandler
	auth *AuthHandler
}

func New(meta *MetaHandler, user *UserHandler, auth *AuthHandler) *Handler {
	return &Handler{
		meta: meta,
		user: user,
		auth: auth,
	}
}

var _ v1Generated.ServerInterface = (*Handler)(nil)

func (h *Handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	h.meta.GetHealth(w, r)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	h.auth.LoginUser(w, r)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.user.RegisterUser(w, r)
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
