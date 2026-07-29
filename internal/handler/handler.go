package handler

import (
	"sketch-api-go/internal/generated"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

var _ generated.ServerInterface = (*Handler)(nil)
