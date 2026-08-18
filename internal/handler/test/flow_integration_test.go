package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sketch-api-go/internal/app"
	v1Generated "sketch-api-go/internal/generated/v1"
	"testing"
)

type testUser struct {
	name         string
	email        string
	password     string
	confirmation string
	accessToken  string
	cookie       *http.Cookie
}

func (u *testUser) registerBody() io.Reader {
	body, _ := json.Marshal(map[string]string{
		"username":     u.name,
		"email":        u.email,
		"password":     u.password,
		"confirmation": u.confirmation,
	})

	return bytes.NewReader(body)
}

func (u *testUser) loginBody() io.Reader {
	body, _ := json.Marshal(map[string]string{
		"email":    u.email,
		"password": u.password,
	})

	return bytes.NewReader(body)
}

func TestAuthFlow(t *testing.T) {

	userTest := testUser{
		name:         "Alice",
		email:        "alice@example.com",
		password:     "secret123",
		confirmation: "secret123",
	}

	fmt.Printf("config = %#v\n", testConfig)
	fmt.Printf("server = %#v\n", testServet)
	fmt.Printf("logger = %#v\n", testLogger)
	fmt.Printf("db = %#v\n", testDB)
	fmt.Printf("jwt = %#v\n", testJWTManager)
	//fmt.Printf("cors = %#v\n", testConfig.CORS())
	fmt.Printf("cookie = %#v\n", testCookieManager)

	serverApp := app.New(
		testServet,
		testLogger,
		testDB,
		testJWTManager,
		testCORS,
		testCookieManager,
	)

	httpHandler := serverApp.Handler()

	// -------------------------------------------------------------------------
	// Endpoint: /user/register

	reqReg := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/user/register",
		userTest.registerBody(),
	)

	reqReg.Header.Set("Content-Type", "application/json")

	recReg := httptest.NewRecorder()

	httpHandler.ServeHTTP(recReg, reqReg)

	if recReg.Code != http.StatusCreated {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			recReg.Code,
			http.StatusCreated,
			recReg.Body.String(),
		)
	}

	var responseReg v1Generated.UserResponse

	if err := json.NewDecoder(recReg.Body).Decode(&responseReg); err != nil {
		t.Fatalf("decode responseReg: %v", err)
	}

	if responseReg.Username != userTest.name {
		t.Fatalf(
			"username = %q, want %q",
			responseReg.Username,
			"Alice",
		)
	}

	// -------------------------------------------------------------------------
	// Endpoint: /auth/login

	reqLogin := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/login",
		userTest.loginBody(),
	)

	reqLogin.Header.Set("Content-Type", "application/json")

	recLogin := httptest.NewRecorder()

	httpHandler.ServeHTTP(recLogin, reqLogin)

	if recLogin.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			recLogin.Code,
			http.StatusOK,
			recLogin.Body.String(),
		)
	}

	var responseLogin v1Generated.TokensResponse

	if err := json.NewDecoder(recLogin.Body).Decode(&responseLogin); err != nil {
		t.Fatalf("decode responseLogin: %v", err)
	}

	if responseLogin.AccessToken == "" {
		t.Fatalf(
			"access_tokin = %q",
			responseLogin.AccessToken,
		)
	}

	userTest.accessToken = responseLogin.AccessToken

	var refreshCookieLogin *http.Cookie

	for _, cookie := range recLogin.Result().Cookies() {
		if cookie.Name == testCookieManager.Name() {
			refreshCookieLogin = cookie
			break
		}
	}

	if refreshCookieLogin == nil || refreshCookieLogin.Value == "" {
		t.Fatal("refresh_token cookie not found")
	}

	userTest.cookie = refreshCookieLogin

	// -------------------------------------------------------------------------
	// Endpoint: /user/me vol.#1

	reqMeV1 := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/user/me",
		nil,
	)

	reqMeV1.Header.Set(
		"Authorization",
		"Bearer "+userTest.accessToken,
	)
	reqMeV1.AddCookie(userTest.cookie)

	recMeV1 := httptest.NewRecorder()

	httpHandler.ServeHTTP(recMeV1, reqMeV1)

	if recMeV1.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			recMeV1.Code,
			http.StatusOK,
			recMeV1.Body.String(),
		)
	}

	var responseMeV1 v1Generated.UserResponse

	if err := json.NewDecoder(recMeV1.Body).Decode(&responseMeV1); err != nil {
		t.Fatalf("decode responseMe: %v", err)
	}

	if responseMeV1.Username != userTest.name {
		t.Fatalf(
			"username = %q, want %q",
			responseMeV1.Username,
			userTest.name,
		)
	}

	// -------------------------------------------------------------------------
	// Endpoint: /auth/refresh vol.#1

	reqRefV1 := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/refresh",
		nil,
	)

	reqRefV1.AddCookie(userTest.cookie)

	recRefV1 := httptest.NewRecorder()

	httpHandler.ServeHTTP(recRefV1, reqRefV1)

	if recRefV1.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			recRefV1.Code,
			http.StatusOK,
			recRefV1.Body.String(),
		)
	}

	var responseRefV1 v1Generated.TokensResponse

	if err := json.NewDecoder(recRefV1.Body).Decode(&responseRefV1); err != nil {
		t.Fatalf("decode responseRefV1: %v", err)
	}

	if responseRefV1.AccessToken == "" {
		t.Fatalf(
			"access_tokin = %q",
			responseRefV1.AccessToken,
		)
	}

	userTest.accessToken = responseRefV1.AccessToken

	var refreshCookieRefV1 *http.Cookie

	for _, cookie := range recRefV1.Result().Cookies() {
		if cookie.Name == testCookieManager.Name() {
			refreshCookieRefV1 = cookie
			break
		}
	}

	if refreshCookieRefV1 == nil || refreshCookieRefV1.Value == "" {
		t.Fatal("refresh_token cookie not found")
	}

	if refreshCookieRefV1.Value == userTest.cookie.Value {
		t.Fatal("refresh token was not rotated")
	}

	userTest.cookie = refreshCookieRefV1

	// -------------------------------------------------------------------------
	// Endpoint: /user/me vol.#2

	reqMeV2 := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/user/me",
		nil,
	)

	reqMeV2.Header.Set(
		"Authorization",
		"Bearer "+userTest.accessToken,
	)
	reqMeV2.AddCookie(userTest.cookie)

	recMeV2 := httptest.NewRecorder()

	httpHandler.ServeHTTP(recMeV2, reqMeV2)

	if recMeV2.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			recMeV2.Code,
			http.StatusOK,
			recMeV2.Body.String(),
		)
	}

	var responseMeV2 v1Generated.UserResponse

	if err := json.NewDecoder(recMeV2.Body).Decode(&responseMeV2); err != nil {
		t.Fatalf("decode responseMe: %v", err)
	}

	if responseMeV2.Username != userTest.name {
		t.Fatalf(
			"username = %q, want %q",
			responseMeV2.Username,
			userTest.name,
		)
	}

	// -------------------------------------------------------------------------
	// Endpoint: /auth/logout

	reqLogout := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/logout",
		nil,
	)

	reqLogout.Header.Set(
		"Authorization",
		"Bearer "+userTest.accessToken,
	)
	reqLogout.AddCookie(userTest.cookie)

	recLogout := httptest.NewRecorder()

	httpHandler.ServeHTTP(recLogout, reqLogout)

	if recLogout.Code != http.StatusNoContent {
		t.Fatalf(
			"status = %d, want %d",
			recLogout.Code,
			http.StatusNoContent,
		)
	}

	var refreshCookieLogout *http.Cookie

	for _, cookie := range recLogout.Result().Cookies() {
		if cookie.Name == testCookieManager.Name() {
			refreshCookieLogout = cookie
			break
		}
	}

	if refreshCookieLogout == nil {
		t.Fatal("refresh_token cookie not found")
	}

	if refreshCookieLogout.Value != "" {
		t.Fatalf(
			"refresh_token cookie value = %q, want empty",
			refreshCookieLogout.Value,
		)
	}

	// -------------------------------------------------------------------------
	// Endpoint: /auth/refresh vol.#2

	reqRefV2 := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/refresh",
		nil,
	)

	reqRefV2.AddCookie(userTest.cookie)

	recRefV2 := httptest.NewRecorder()

	httpHandler.ServeHTTP(recRefV2, reqRefV2)

	if recRefV2.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want %d; body = %s",
			recRefV2.Code,
			http.StatusUnauthorized,
			recRefV2.Body.String(),
		)
	}
}
