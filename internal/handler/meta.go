package handler

import (
	"net/http"
)

type MetaHandler struct{}

func (m *MetaHandler) GetHealth(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("OK"))
}
