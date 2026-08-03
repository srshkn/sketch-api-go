package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sketch-api-go/internal/db"
	"sketch-api-go/internal/generated"
	"sketch-api-go/internal/repository"
	"sketch-api-go/internal/service"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type userRepositoryStub struct {
	repository.UserRepository

	createUserFn func(
		context.Context,
		db.CreateUserParams,
	) (db.CreateUserRow, error)
}

func (s userRepositoryStub) CreateUser(
	ctx context.Context,
	params db.CreateUserParams,
) (db.CreateUserRow, error) {
	return s.createUserFn(ctx, params)
}

func newUserTestRouter(repo repository.UserRepository) http.Handler {
	userService := service.NewUserService(repo)
	userHandler := NewUserHandler(userService)
	handler := New(NewMetaHandler(), userHandler)

	return generated.HandlerFromMux(handler,
		http.NewServeMux())
}

func TestRegisterUser(t *testing.T) {
	wantID := uuid.MustParse("152c35c2-167c-47c9-891a-f8caf6474eaf")

	repo := userRepositoryStub{
		createUserFn: func(
			ctx context.Context,
			params db.CreateUserParams,
		) (db.CreateUserRow, error) {
			if params.Name != "Alice" {
				t.Errorf("name = %q, want %q", params.Name, "Alice")
			}

			if params.Email != "alice@example.com" {
				t.Errorf(
					"email = %q, want %q",
					params.Email,
					"alice@example.com",
				)
			}

			return db.CreateUserRow{
				ID:    wantID,
				Name:  params.Name,
				Email: params.Email,
			}, nil
		},
	}

	router := newUserTestRouter(repo)

	request := httptest.NewRequest(
		http.MethodPost,
		"/user/register",
		strings.NewReader(`{
              "username": "  Alice  ",
              "email": "alice@example.com",
              "password": "secret123",
              "confirmation": "secret123"
          }`),
	)

	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf(
			"status code = %d, want %d; body = %s",
			response.Code,
			http.StatusCreated,
			response.Body.String(),
		)
	}

	var body generated.UserResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if body.Id != wantID {
		t.Errorf("id = %s, want %s", body.Id, wantID)
	}

	if body.Username != "Alice" {
		t.Errorf("username = %q, want %q", body.Username, "Alice")
	}
}

/*
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
*/
