package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"backend/models"

	beego "github.com/beego/beego/v2/server/web"
)

type expenseResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    expenseResponseData `json:"data"`
}

type expenseResponseData struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

func TestCreateExpenseSuccessReturnsCreated(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := postJSONWithHeaders("/api/v1/expenses", validExpenseRequestBody(), map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if !response.Success {
		t.Fatal("expected success to be true")
	}
	if response.Message != "Expense created successfully" {
		t.Fatalf("expected message %q, got %q", "Expense created successfully", response.Message)
	}
	if response.Data.ID != 1 || response.Data.Title != "Lunch" || response.Data.Amount != 350.50 || response.Data.Category != "Food" || response.Data.Note != "Team lunch" || response.Data.ExpenseDate != "2025-06-10" {
		t.Fatalf("unexpected expense response data: %+v", response.Data)
	}
}

func TestCreateExpenseStoresAuthenticatedUserID(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := postJSONWithHeaders("/api/v1/expenses", validExpenseRequestBody(), map[string]string{"X-User-ID": "1"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	expenses, err := models.GetExpensesByUserID(1)
	if err != nil {
		t.Fatalf("expected expenses to load: %v", err)
	}
	if len(expenses) != 1 {
		t.Fatalf("expected one expense, got %d", len(expenses))
	}
	if expenses[0].UserID != 1 {
		t.Fatalf("expected authenticated user ID 1, got %d", expenses[0].UserID)
	}
}

func TestCreateExpenseMissingUserIDReturnsUnauthorized(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := postJSON("/api/v1/expenses", validExpenseRequestBody())
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusUnauthorized, "Unauthorized")
}

func TestCreateExpenseInvalidUserIDReturnsUnauthorized(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := postJSONWithHeaders("/api/v1/expenses", validExpenseRequestBody(), map[string]string{"X-User-ID": "invalid"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusUnauthorized, "Unauthorized")
}

func TestCreateExpenseNonExistingUserIDReturnsUnauthorized(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := postJSONWithHeaders("/api/v1/expenses", validExpenseRequestBody(), map[string]string{"X-User-ID": "99"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusUnauthorized, "Unauthorized")
}

func TestCreateExpenseInvalidJSONReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := postJSONWithHeaders("/api/v1/expenses", `{"title":"Lunch"`, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid request body")
}

func TestCreateExpenseMissingTitleReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	body := `{"amount":350.50,"category":"Food","note":"Team lunch","expense_date":"2025-06-10"}`
	rec := postJSONWithHeaders("/api/v1/expenses", body, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Title is required")
}

func TestCreateExpenseZeroAmountReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	body := `{"title":"Lunch","amount":0,"category":"Food","note":"Team lunch","expense_date":"2025-06-10"}`
	rec := postJSONWithHeaders("/api/v1/expenses", body, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Amount must be positive")
}

func TestCreateExpenseInvalidCategoryReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	body := `{"title":"Lunch","amount":350.50,"category":"Travel","note":"Team lunch","expense_date":"2025-06-10"}`
	rec := postJSONWithHeaders("/api/v1/expenses", body, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid category")
}

func TestCreateExpenseInvalidExpenseDateReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	body := `{"title":"Lunch","amount":350.50,"category":"Food","note":"Team lunch","expense_date":"2025/06/10"}`
	rec := postJSONWithHeaders("/api/v1/expenses", body, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid expense date format")
}

func useTempUserAndExpenseCSVConfig(t *testing.T) {
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

func seedAuthenticatedUser(t *testing.T, userID int) {
	t.Helper()

	err := models.CreateUser(&models.User{
		ID:       userID,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected authenticated user to be seeded: %v", err)
	}
}

func postJSONWithHeaders(path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(rec, req)

	return rec
}

func decodeExpenseResponse(t *testing.T, rec *httptest.ResponseRecorder) expenseResponse {
	t.Helper()

	var response expenseResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return response
}

func assertExpenseErrorResponse(t *testing.T, statusCode int, response expenseResponse, expectedStatus int, expectedMessage string) {
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

func validExpenseRequestBody() string {
	return `{"title":"Lunch","amount":350.50,"category":"Food","note":"Team lunch","expense_date":"2025-06-10"}`
}
