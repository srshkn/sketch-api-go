package v1

import (
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

// -------------------------------------------------------------------------
// Meta

func (h *Handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	h.meta.GetHealth(w, r)
}

// -------------------------------------------------------------------------
// User

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.user.RegisterUser(w, r)
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	h.user.GetMe(w, r)
}

// -------------------------------------------------------------------------
// Auth

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	h.auth.LoginUser(w, r)
}

func (h *Handler) UpdateRefreshToken(w http.ResponseWriter, r *http.Request) {
	h.auth.UpdateRefreshToken(w, r)
}
