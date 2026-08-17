package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	v1Generated "sketch-api-go/internal/generated/v1"
	v1 "sketch-api-go/internal/handler/v1"
	"sketch-api-go/internal/service"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestRegisterUserIntegration(t *testing.T) {
	svc := service.NewUserService(testDB)

	handler := v1.NewUserHandler(svc)

	router := v1.New(
		v1.NewMetaHandler(),
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
			"username": "Ben",
			"email": "benben@example.com",
			"password": "secret321",
			"confirmation": "secret321"
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

	if response.Username != "Ben" {
		t.Errorf(
			"username = %q, want %q",
			response.Username,
			"Ben",
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

	if user.Email != "benben@example.com" {
		t.Errorf(
			"email = %q, want %q",
			user.Email,
			"benben@example.com",
		)
	}

	if user.Username != "Ben" {
		t.Errorf(
			"database username = %q, want %q",
			user.Username,
			"Ben",
		)
	}

}
