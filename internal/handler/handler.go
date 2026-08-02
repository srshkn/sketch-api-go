package handler

import (
	"net/http"
	"sketch-api-go/internal/generated"
)

type Handler struct {
	meta *MetaHandler
	user *UserHandler
}

func New() *Handler {
	return &Handler{
		meta: &MetaHandler{},
		user: &UserHandler{},
	}
}

var _ generated.ServerInterface = (*Handler)(nil)

func (h *Handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	h.meta.GetHealth(w, r)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.user.RegisterUser(w, r)
}
