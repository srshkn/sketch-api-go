package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"sketch-api-go/internal/service"

	v1Generated "sketch-api-go/internal/generated/v1"
)

type userServiceStub struct {
	registrationFn func(
		context.Context,
		v1Generated.RegisterUserRequest,
	) (v1Generated.UserResponse, error)

	getUserFn func(
		context.Context,
		string,
	) (v1Generated.UserResponse, error)
}

func (s userServiceStub) Registration(
	ctx context.Context,
	request v1Generated.RegisterUserRequest,
) (v1Generated.UserResponse, error) {
	return s.registrationFn(ctx, request)
}

func (s userServiceStub) GetUser(
	ctx context.Context,
	userID string,
) (v1Generated.UserResponse, error) {
	return s.getUserFn(ctx, userID)
}

func newUserTestRouter(userService service.User) http.Handler {
	userHandler := NewUserHandler(userService)
	handler := New(NewMetaHandler(), userHandler, nil)

	return v1Generated.HandlerFromMux(
		handler,
		http.NewServeMux(),
	)
}

func TestRegisterUser(t *testing.T) {
	wantID := uuid.MustParse("152c35c2-167c-47c9-891a-f8caf6474eaf")

	serviceStub := userServiceStub{
		registrationFn: func(
			ctx context.Context,
			request v1Generated.RegisterUserRequest,
		) (v1Generated.UserResponse, error) {
			if request.Username != "Alice" {
				t.Errorf(
					"username = %q, want %q",
					request.Username,
					"Alice",
				)
			}

			if request.Password != "secret123" {
				t.Errorf(
					"password = %q, want %q",
					request.Password,
					"secret123",
				)
			}

			return v1Generated.UserResponse{
				Id:       wantID,
				Username: request.Username,
			}, nil
		},
	}

	router := newUserTestRouter(serviceStub)

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

	var body v1Generated.UserResponse

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if body.Id != wantID {
		t.Errorf(
			"id = %s, want %s",
			body.Id,
			wantID,
		)
	}

	if body.Username != "Alice" {
		t.Errorf(
			"username = %q, want %q",
			body.Username,
			"Alice",
		)
	}
}

func TestRegisterUserValidation(t *testing.T) {
	serviceCalled := false

	serviceStub := userServiceStub{
		registrationFn: func(
			ctx context.Context,
			request v1Generated.RegisterUserRequest,
		) (v1Generated.UserResponse, error) {
			serviceCalled = true
			return v1Generated.UserResponse{}, nil
		},
	}

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
			body:        `{"username":`,
			wantMessage: "invalid request body",
		},
		{
			name:        "missing username",
			body:        `{"password":"secret"}`,
			wantMessage: "username must not be empty",
		},
		{
			name:        "blank username",
			body:        `{"username":"   ","password":"secret"}`,
			wantMessage: "username must not be empty",
		},
		{
			name:        "missing password",
			body:        `{"username":"Alice"}`,
			wantMessage: "password must not be empty",
		},
		{
			name:        "empty password",
			body:        `{"username":"Alice","password":""}`,
			wantMessage: "password must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceCalled = false

			router := newUserTestRouter(serviceStub)

			request := httptest.NewRequest(
				http.MethodPost,
				"/user/register",
				strings.NewReader(tt.body),
			)

			request.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusBadRequest {
				t.Fatalf(
					"status code = %d, want %d",
					response.Code,
					http.StatusBadRequest,
				)
			}

			var body v1Generated.ErrorResponse

			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatalf(
					"decode response body: %v",
					err,
				)
			}

			if body.Error.Code != v1Generated.INVALIDREQUEST {
				t.Errorf(
					"error code = %q, want %q",
					body.Error.Code,
					v1Generated.INVALIDREQUEST,
				)
			}

			if body.Error.Message != tt.wantMessage {
				t.Errorf(
					"error message = %q, want %q",
					body.Error.Message,
					tt.wantMessage,
				)
			}

			if serviceCalled {
				t.Error("service was called for invalid request")
			}
		})
	}
}

func TestRegisterUserServiceError(t *testing.T) {
	serviceStub := userServiceStub{
		registrationFn: func(
			ctx context.Context,
			request v1Generated.RegisterUserRequest,
		) (v1Generated.UserResponse, error) {
			return v1Generated.UserResponse{}, errors.New("database error")
		},
	}

	router := newUserTestRouter(serviceStub)

	request := httptest.NewRequest(
		http.MethodPost,
		"/user/register",
		strings.NewReader(`{
			"username": "Alice",
			"email": "alice@example.com",
			"password": "secret123",
			"confirmation": "secret123"
		}`),
	)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf(
			"status code = %d, want %d",
			response.Code,
			http.StatusBadRequest,
		)
	}
}
