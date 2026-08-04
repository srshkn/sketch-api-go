package handler

import (
	"net/http"

	v1Generated "sketch-api-go/internal/generated/v1"
)

type Handler struct {
	meta *MetaHandler
	user *UserHandler
}

func New(meta *MetaHandler, user *UserHandler) *Handler {
	return &Handler{
		meta: meta,
		user: user,
	}
}

var _ v1Generated.ServerInterface = (*Handler)(nil)

func (h *Handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	h.meta.GetHealth(w, r)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.user.RegisterUser(w, r)
}
