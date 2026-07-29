package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sketch-api-go/internal/generated"
)

func TestGetHealth(t *testing.T) {
	router := generated.HandlerFromMux(New(), http.NewServeMux())
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

func TestRegisterUser(t *testing.T) {
	router := generated.HandlerFromMux(New(), http.NewServeMux())
	request := httptest.NewRequest(
		http.MethodPost,
		"/user/register",
		strings.NewReader(`{"name":"  Alice  ","password":"secret"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusCreated)
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	var body generated.UserResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if body.Id != 1 {
		t.Errorf("id = %d, want 1", body.Id)
	}
	if body.Name != "Alice" {
		t.Errorf("name = %q, want %q", body.Name, "Alice")
	}
}

func TestRegisterUserValidation(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantMessage string
	}{
		{
			name:        "empty body",
			body:        "",
			wantMessage: "invalid request body",
		},
		{
			name:        "malformed JSON",
			body:        `{"name":`,
			wantMessage: "invalid request body",
		},
		{
			name:        "unknown field",
			body:        `{"name":"Alice","password":"secret","role":"admin"}`,
			wantMessage: "invalid request body",
		},
		{
			name:        "missing name",
			body:        `{"password":"secret"}`,
			wantMessage: "name must not be empty",
		},
		{
			name:        "blank name",
			body:        `{"name":"   ","password":"secret"}`,
			wantMessage: "name must not be empty",
		},
		{
			name:        "missing password",
			body:        `{"name":"Alice"}`,
			wantMessage: "password must not be empty",
		},
		{
			name:        "empty password",
			body:        `{"name":"Alice","password":""}`,
			wantMessage: "password must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := generated.HandlerFromMux(New(), http.NewServeMux())
			request := httptest.NewRequest(
				http.MethodPost,
				"/user/register",
				strings.NewReader(tt.body),
			)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf("status code = %d, want %d", response.Code, http.StatusBadRequest)
			}
			if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
			}

			var body generated.ErrorResponse
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatalf("decode response body: %v", err)
			}
			if body.Error.Code != generated.INVALIDREQUEST {
				t.Errorf("error code = %q, want %q", body.Error.Code, generated.INVALIDREQUEST)
			}
			if body.Error.Message != tt.wantMessage {
				t.Errorf("error message = %q, want %q", body.Error.Message, tt.wantMessage)
			}
		})
	}
}
