package routers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"backend/models"
	_ "backend/routers"

	beego "github.com/beego/beego/v2/server/web"
)

type routerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	appPath := filepath.Dir(filepath.Dir(file))
	beego.TestBeegoInit(appPath)
}

func TestHealthRouteIsRegistered(t *testing.T) {
	rec := serveRouterRequest(http.MethodGet, "/api/v1/health", "", nil)
	response := decodeRouterResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !response.Success {
		t.Fatal("expected success to be true")
	}
	if response.Message != "Server is running" {
		t.Fatalf("expected message %q, got %q", "Server is running", response.Message)
	}
}

func TestSummaryRouteIsRegisteredBeforeExpenseIDRoute(t *testing.T) {
	useTempCSVConfig(t)
	seedRouterUser(t, 1)

	rec := serveRouterRequest(
		http.MethodGet,
		"/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30",
		"",
		map[string]string{"X-User-ID": "1"},
	)
	response := decodeRouterResponse(t, rec)

	if response.Message == "Invalid expense ID" {
		t.Fatal("expected summary route, got expense ID route")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %q", http.StatusOK, rec.Code, rec.Body.String())
	}
	if !response.Success {
		t.Fatal("expected success to be true")
	}
	if response.Message != "Summary generated" {
		t.Fatalf("expected message %q, got %q", "Summary generated", response.Message)
	}
}

func TestMissingExpenseIDRouteIsRegistered(t *testing.T) {
	useTempCSVConfig(t)
	seedRouterUser(t, 1)

	rec := serveRouterRequest(
		http.MethodGet,
		"/api/v1/expenses/999",
		"",
		map[string]string{"X-User-ID": "1"},
	)
	response := decodeRouterResponse(t, rec)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if response.Success {
		t.Fatal("expected success to be false")
	}
	if response.Message != "Expense not found" {
		t.Fatalf("expected message %q, got %q", "Expense not found", response.Message)
	}
}

func TestImportantMethodsAreRegistered(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		headers    map[string]string
		wantStatus int
	}{
		{
			name:       "register post route",
			method:     http.MethodPost,
			path:       "/api/v1/auth/register",
			body:       `{"name":"John Doe"`,
			headers:    map[string]string{"Content-Type": "application/json"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "login post route",
			method:     http.MethodPost,
			path:       "/api/v1/auth/login",
			body:       `{"email":"john@example.com"`,
			headers:    map[string]string{"Content-Type": "application/json"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "create expense post route",
			method:     http.MethodPost,
			path:       "/api/v1/expenses",
			body:       `{}`,
			headers:    map[string]string{"Content-Type": "application/json"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "list expense get route",
			method:     http.MethodGet,
			path:       "/api/v1/expenses",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "update expense put route",
			method:     http.MethodPut,
			path:       "/api/v1/expenses/999",
			body:       `{}`,
			headers:    map[string]string{"Content-Type": "application/json"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "delete expense delete route",
			method:     http.MethodDelete,
			path:       "/api/v1/expenses/999",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serveRouterRequest(tt.method, tt.path, tt.body, tt.headers)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d with body %q", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func useTempCSVConfig(t *testing.T) {
	t.Helper()

	previousUserPath := beego.AppConfig.DefaultString("csv_user_file", "data/users.csv")
	previousExpensePath := beego.AppConfig.DefaultString("csv_expense_file", "data/expenses.csv")
	tempDir := t.TempDir()

	if err := beego.AppConfig.Set("csv_user_file", filepath.Join(tempDir, "users.csv")); err != nil {
		t.Fatalf("expected csv_user_file test config to be set: %v", err)
	}
	if err := beego.AppConfig.Set("csv_expense_file", filepath.Join(tempDir, "expenses.csv")); err != nil {
		t.Fatalf("expected csv_expense_file test config to be set: %v", err)
	}

	t.Cleanup(func() {
		_ = beego.AppConfig.Set("csv_user_file", previousUserPath)
		_ = beego.AppConfig.Set("csv_expense_file", previousExpensePath)
	})
}

func seedRouterUser(t *testing.T, userID int) {
	t.Helper()

	if err := models.CreateUser(&models.User{
		ID:       userID,
		Name:     "Router User",
		Email:    "router@example.com",
		Password: "secret123",
	}); err != nil {
		t.Fatalf("expected router user to be seeded: %v", err)
	}
}

func serveRouterRequest(method string, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(rec, req)

	return rec
}

func decodeRouterResponse(t *testing.T, rec *httptest.ResponseRecorder) routerResponse {
	t.Helper()

	var response routerResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return response
}
