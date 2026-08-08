package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	v1Generated "sketch-api-go/internal/generated/v1"
)

func TestGetHealth(t *testing.T) {
	handlerMeta := New(NewMetaHandler(), nil, nil)
	router := v1Generated.HandlerFromMux(handlerMeta, http.NewServeMux())
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "text/plain; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", contentType, "text/plain; charset=utf-8")
	}
	if body := response.Body.String(); body != "OK" {
		t.Errorf("body = %q, want %q", body, "OK")
	}
}
