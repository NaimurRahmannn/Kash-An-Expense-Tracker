package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"backend/models"
	"backend/utils"

	beego "github.com/beego/beego/v2/server/web"
)

type authResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    authResponseData `json:"data"`
}

type authResponseData struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func TestRegisterSuccessReturnsCreated(t *testing.T) {
	useTempUserCSVConfig(t)

	rec := postJSON("/api/v1/auth/register", `{"name":"John Doe","email":"john@example.com","password":"secret123"}`)
	response := decodeAuthResponse(t, rec)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if !response.Success {
		t.Fatal("expected success to be true")
	}
	if response.Message != "User registered successfully" {
		t.Fatalf("expected message %q, got %q", "User registered successfully", response.Message)
	}
	if strings.Contains(rec.Body.String(), "password") {
		t.Fatal("expected register response to omit password")
	}
}

func TestRegisterStoresHashedPassword(t *testing.T) {
	useTempUserCSVConfig(t)

	postJSON("/api/v1/auth/register", `{"name":"John Doe","email":"john@example.com","password":"secret123"}`)

	user, err := models.GetUserByEmail("john@example.com")
	if err != nil {
		t.Fatalf("expected user lookup to succeed: %v", err)
	}
	if user == nil {
		t.Fatal("expected user to be stored")
	}
	if user.Password == "secret123" {
		t.Fatal("expected stored password to be hashed")
	}
	if !utils.CheckPasswordHash("secret123", user.Password) {
		t.Fatal("expected stored password to match bcrypt hash")
	}
}

func TestRegisterDuplicateEmailReturnsConflict(t *testing.T) {
	useTempUserCSVConfig(t)

	postJSON("/api/v1/auth/register", `{"name":"John Doe","email":"john@example.com","password":"secret123"}`)
	rec := postJSON("/api/v1/auth/register", `{"name":"Jane Doe","email":"john@example.com","password":"secret456"}`)
	response := decodeAuthResponse(t, rec)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
	if response.Success {
		t.Fatal("expected success to be false")
	}
	if response.Message != "Email already exists" {
		t.Fatalf("expected message %q, got %q", "Email already exists", response.Message)
	}
}

func TestRegisterMissingNameReturnsBadRequest(t *testing.T) {
	useTempUserCSVConfig(t)

	rec := postJSON("/api/v1/auth/register", `{"email":"john@example.com","password":"secret123"}`)
	response := decodeAuthResponse(t, rec)

	assertErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Name is required")
}

func TestRegisterInvalidEmailReturnsBadRequest(t *testing.T) {
	useTempUserCSVConfig(t)

	rec := postJSON("/api/v1/auth/register", `{"name":"John Doe","email":"invalid-email","password":"secret123"}`)
	response := decodeAuthResponse(t, rec)

	assertErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid email format")
}

func TestRegisterShortPasswordReturnsBadRequest(t *testing.T) {
	useTempUserCSVConfig(t)

	rec := postJSON("/api/v1/auth/register", `{"name":"John Doe","email":"john@example.com","password":"12345"}`)
	response := decodeAuthResponse(t, rec)

	assertErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Password must be at least 6 characters")
}

func TestRegisterValidationFailuresReturnBadRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantMessage string
	}{
		{
			name:        "missing email",
			body:        `{"name":"John Doe","password":"secret123"}`,
			wantMessage: "Email is required",
		},
		{
			name:        "missing password",
			body:        `{"name":"John Doe","email":"john@example.com"}`,
			wantMessage: "Password is required",
		},
		{
			name:        "invalid JSON",
			body:        `{"name":"John Doe"`,
			wantMessage: "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempUserCSVConfig(t)

			rec := postJSON("/api/v1/auth/register", tt.body)
			response := decodeAuthResponse(t, rec)

			assertErrorResponse(t, rec.Code, response, http.StatusBadRequest, tt.wantMessage)
		})
	}
}

func TestRegisterReturnsInternalServerErrorWhenUserLookupFails(t *testing.T) {
	useMalformedUserCSVConfig(t)

	rec := postJSON("/api/v1/auth/register", `{"name":"John Doe","email":"john@example.com","password":"secret123"}`)
	response := decodeAuthResponse(t, rec)

	assertErrorResponse(t, rec.Code, response, http.StatusInternalServerError, "Internal server error")
}

func TestLoginSuccessReturnsUserData(t *testing.T) {
	useTempUserCSVConfig(t)

	postJSON("/api/v1/auth/register", `{"name":"John Doe","email":"john@example.com","password":"secret123"}`)
	rec := postJSON("/api/v1/auth/login", `{"email":"john@example.com","password":"secret123"}`)
	response := decodeAuthResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !response.Success {
		t.Fatal("expected success to be true")
	}
	if response.Message != "Login successful" {
		t.Fatalf("expected message %q, got %q", "Login successful", response.Message)
	}
	if response.Data.UserID != 1 {
		t.Fatalf("expected user ID 1, got %d", response.Data.UserID)
	}
	if response.Data.Name != "John Doe" {
		t.Fatalf("expected name %q, got %q", "John Doe", response.Data.Name)
	}
	if response.Data.Email != "john@example.com" {
		t.Fatalf("expected email %q, got %q", "john@example.com", response.Data.Email)
	}
	if strings.Contains(rec.Body.String(), "secret123") || strings.Contains(rec.Body.String(), "password") {
		t.Fatal("expected login response to omit password")
	}
}

func TestLoginWrongPasswordReturnsUnauthorized(t *testing.T) {
	useTempUserCSVConfig(t)

	postJSON("/api/v1/auth/register", `{"name":"John Doe","email":"john@example.com","password":"secret123"}`)
	rec := postJSON("/api/v1/auth/login", `{"email":"john@example.com","password":"wrongpassword"}`)
	response := decodeAuthResponse(t, rec)

	assertErrorResponse(t, rec.Code, response, http.StatusUnauthorized, "Invalid email or password")
}

func TestLoginPlaintextPasswordReturnsUnauthorized(t *testing.T) {
	useTempUserCSVConfig(t)

	if err := models.CreateUser(&models.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}); err != nil {
		t.Fatalf("expected user to be seeded: %v", err)
	}

	rec := postJSON("/api/v1/auth/login", `{"email":"john@example.com","password":"secret123"}`)
	response := decodeAuthResponse(t, rec)

	assertErrorResponse(t, rec.Code, response, http.StatusUnauthorized, "Invalid email or password")
}

func TestLoginMissingEmailReturnsBadRequest(t *testing.T) {
	useTempUserCSVConfig(t)

	rec := postJSON("/api/v1/auth/login", `{"password":"secret123"}`)
	response := decodeAuthResponse(t, rec)

	assertErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Email is required")
}

func TestAuthInvalidJSONReturnsBadRequest(t *testing.T) {
	useTempUserCSVConfig(t)

	rec := postJSON("/api/v1/auth/register", `{"name":"John Doe"`)
	response := decodeAuthResponse(t, rec)

	assertErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid request body")
}

func TestLoginValidationFailuresReturnBadRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantMessage string
	}{
		{
			name:        "invalid email",
			body:        `{"email":"invalid-email","password":"secret123"}`,
			wantMessage: "Invalid email format",
		},
		{
			name:        "missing password",
			body:        `{"email":"john@example.com"}`,
			wantMessage: "Password is required",
		},
		{
			name:        "invalid JSON",
			body:        `{"email":"john@example.com"`,
			wantMessage: "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useTempUserCSVConfig(t)

			rec := postJSON("/api/v1/auth/login", tt.body)
			response := decodeAuthResponse(t, rec)

			assertErrorResponse(t, rec.Code, response, http.StatusBadRequest, tt.wantMessage)
		})
	}
}

func TestLoginReturnsInternalServerErrorWhenUserLookupFails(t *testing.T) {
	useMalformedUserCSVConfig(t)

	rec := postJSON("/api/v1/auth/login", `{"email":"john@example.com","password":"secret123"}`)
	response := decodeAuthResponse(t, rec)

	assertErrorResponse(t, rec.Code, response, http.StatusInternalServerError, "Internal server error")
}

func useTempUserCSVConfig(t *testing.T) {
	t.Helper()

	previousPath := beego.AppConfig.DefaultString("csv_user_file", "data/users.csv")
	filePath := filepath.Join(t.TempDir(), "users.csv")

	if err := beego.AppConfig.Set("csv_user_file", filePath); err != nil {
		t.Fatalf("expected csv_user_file test config to be set: %v", err)
	}

	t.Cleanup(func() {
		_ = beego.AppConfig.Set("csv_user_file", previousPath)
	})
}

func useMalformedUserCSVConfig(t *testing.T) {
	t.Helper()

	previousPath := beego.AppConfig.DefaultString("csv_user_file", "data/users.csv")
	filePath := filepath.Join(t.TempDir(), "users.csv")
	content := "id,name,email,password,created_at\nbad-id,John Doe,john@example.com,secret123,2025-06-01T10:30:00Z\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("expected malformed users CSV to be written: %v", err)
	}
	if err := beego.AppConfig.Set("csv_user_file", filePath); err != nil {
		t.Fatalf("expected csv_user_file test config to be set: %v", err)
	}

	t.Cleanup(func() {
		_ = beego.AppConfig.Set("csv_user_file", previousPath)
	})
}

func postJSON(path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(rec, req)

	return rec
}

func decodeAuthResponse(t *testing.T, rec *httptest.ResponseRecorder) authResponse {
	t.Helper()

	var response authResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return response
}

func assertErrorResponse(t *testing.T, statusCode int, response authResponse, expectedStatus int, expectedMessage string) {
	t.Helper()

	if statusCode != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, statusCode)
	}
	if response.Success {
		t.Fatal("expected success to be false")
	}
	if response.Message != expectedMessage {
		t.Fatalf("expected message %q, got %q", expectedMessage, response.Message)
	}
}
