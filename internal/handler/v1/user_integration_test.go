package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	v1Generated "sketch-api-go/internal/generated/v1"
	"sketch-api-go/internal/service"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRegisterUserIntegration(t *testing.T) {
	svc := service.NewUserService(testDB)

	handler := NewUserHandler(svc)

	router := New(
		NewMetaHandler(),
		handler,
		nil,
	)

	httpHandler := v1Generated.HandlerFromMux(
		router,
		http.NewServeMux(),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/user/register",
		strings.NewReader(`{
			"username": "Alice",
			"email": "alice@example.com",
			"password": "secret123",
			"confirmation": "secret123"
		}`),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	httpHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			rec.Code,
			http.StatusCreated,
			rec.Body.String(),
		)
	}

	var response v1Generated.UserResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Username != "Alice" {
		t.Errorf(
			"username = %q, want %q",
			response.Username,
			"Alice",
		)
	}

	// 10. Проверяем БД.
	ctx := t.Context()

	user, err := testDB.GetUserByID(ctx, response.Id)
	if err != nil {
		t.Fatalf("query user: %v", err)
	}

	if response.Id == uuid.Nil {
		t.Fatal("response id is empty")
	}

	if user.Email != "alice@example.com" {
		t.Errorf(
			"email = %q, want %q",
			user.Email,
			"alice@example.com",
		)
	}

	if user.Username != "Alice" {
		t.Errorf(
			"database username = %q, want %q",
			user.Username,
			"Alice",
		)
	}

}
