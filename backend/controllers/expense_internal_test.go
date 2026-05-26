package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"backend/utils"

	beego "github.com/beego/beego/v2/server/web"
)

var registerExpenseInternalTestRoutesOnce sync.Once

func TestHandleExpenseWriteErrorResponses(t *testing.T) {
	registerExpenseInternalTestRoutes()

	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "not found error",
			path:        "/test/expense-write-not-found",
			wantStatus:  http.StatusNotFound,
			wantMessage: "Expense not found",
		},
		{
			name:        "internal error",
			path:        "/test/expense-write-internal",
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "Failed to update expense",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serveInternalExpenseRequest(http.MethodGet, tt.path)
			response := decodeInternalAPIResponse(t, rec)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
			if response.Success {
				t.Fatal("expected success to be false")
			}
			if response.Message != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, response.Message)
			}
		})
	}
}

func TestGetExistingExpenseInternalErrorResponse(t *testing.T) {
	registerExpenseInternalTestRoutes()
	useMalformedExpenseCSVConfig(t)

	rec := serveInternalExpenseRequest(http.MethodGet, "/test/existing-expense-internal")
	response := decodeInternalAPIResponse(t, rec)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if response.Success {
		t.Fatal("expected success to be false")
	}
	if response.Message != "Failed to get expense" {
		t.Fatalf("expected message %q, got %q", "Failed to get expense", response.Message)
	}
}

func registerExpenseInternalTestRoutes() {
	registerExpenseInternalTestRoutesOnce.Do(func() {
		beego.Router("/test/expense-write-not-found", &testExpenseWriteNotFoundController{})
		beego.Router("/test/expense-write-internal", &testExpenseWriteInternalController{})
		beego.Router("/test/existing-expense-internal", &testExistingExpenseInternalController{})
	})
}

type testExpenseWriteNotFoundController struct {
	ExpenseController
}

func (c *testExpenseWriteNotFoundController) Get() {
	c.handleExpenseWriteError(errors.New("expense not found"), "Failed to update expense")
}

type testExpenseWriteInternalController struct {
	ExpenseController
}

func (c *testExpenseWriteInternalController) Get() {
	c.handleExpenseWriteError(errors.New("disk write failed"), "Failed to update expense")
}

type testExistingExpenseInternalController struct {
	ExpenseController
}

func (c *testExistingExpenseInternalController) Get() {
	c.getExistingExpense(1, 1, "Failed to get expense")
}

func serveInternalExpenseRequest(method string, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(rec, req)

	return rec
}

func decodeInternalAPIResponse(t *testing.T, rec *httptest.ResponseRecorder) apiResponse {
	t.Helper()

	var response apiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	return response
}

func useMalformedExpenseCSVConfig(t *testing.T) {
	t.Helper()

	previousPath := beego.AppConfig.DefaultString("csv_expense_file", "data/expenses.csv")
	filePath := filepath.Join(t.TempDir(), "expenses.csv")
	rows := [][]string{
		{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"},
		{"1", "1", "Lunch", "bad-amount", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	}

	if err := utils.WriteCSV(filePath, rows); err != nil {
		t.Fatalf("expected malformed expenses CSV to be written: %v", err)
	}
	if err := beego.AppConfig.Set("csv_expense_file", filePath); err != nil {
		t.Fatalf("expected csv_expense_file config to be set: %v", err)
	}

	t.Cleanup(func() {
		_ = beego.AppConfig.Set("csv_expense_file", previousPath)
	})
}
