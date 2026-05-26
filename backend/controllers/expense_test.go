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

type expenseListResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    []expenseResponseData `json:"data"`
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

func TestListExpensesSuccessReturnsOK(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := getWithHeaders("/api/v1/expenses", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseListResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !response.Success {
		t.Fatal("expected success to be true")
	}
	if response.Message != "Expenses retrieved" {
		t.Fatalf("expected message %q, got %q", "Expenses retrieved", response.Message)
	}
	if len(response.Data) != 1 {
		t.Fatalf("expected one expense, got %d", len(response.Data))
	}
}

func TestListExpensesReturnsOnlyAuthenticatedUsersExpenses(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedAuthenticatedUser(t, 2)
	seedExpense(t, 1, "Lunch")
	seedExpense(t, 2, "Bus")

	rec := getWithHeaders("/api/v1/expenses", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseListResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if len(response.Data) != 1 {
		t.Fatalf("expected one expense, got %d", len(response.Data))
	}
	if response.Data[0].Title != "Lunch" {
		t.Fatalf("expected authenticated user's expense, got %+v", response.Data[0])
	}
}

func TestListExpensesWithNoExpensesReturnsEmptyArray(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := getWithHeaders("/api/v1/expenses", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseListResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if len(response.Data) != 0 {
		t.Fatalf("expected empty expenses, got %d", len(response.Data))
	}
}

func TestListExpensesDefaultPaginationWorks(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	for i := 0; i < 11; i++ {
		seedExpense(t, 1, "Lunch")
	}

	rec := getWithHeaders("/api/v1/expenses", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseListResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if len(response.Data) != 10 {
		t.Fatalf("expected default limit of 10 expenses, got %d", len(response.Data))
	}
}

func TestListExpensesPageOneLimitOneReturnsFirstItem(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")
	seedExpense(t, 1, "Dinner")

	rec := getWithHeaders("/api/v1/expenses?page=1&limit=1", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseListResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if len(response.Data) != 1 || response.Data[0].Title != "Lunch" {
		t.Fatalf("expected first expense, got %+v", response.Data)
	}
}

func TestListExpensesPageTwoLimitOneReturnsSecondItem(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")
	seedExpense(t, 1, "Dinner")

	rec := getWithHeaders("/api/v1/expenses?page=2&limit=1", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseListResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if len(response.Data) != 1 || response.Data[0].Title != "Dinner" {
		t.Fatalf("expected second expense, got %+v", response.Data)
	}
}

func TestListExpensesBeyondAvailableDataReturnsEmptyArray(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := getWithHeaders("/api/v1/expenses?page=2&limit=10", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseListResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if len(response.Data) != 0 {
		t.Fatalf("expected empty expenses, got %d", len(response.Data))
	}
}

func TestListExpensesInvalidPageReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := getWithHeaders("/api/v1/expenses?page=0", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid page parameter")
}

func TestListExpensesInvalidLimitReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := getWithHeaders("/api/v1/expenses?limit=abc", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid limit parameter")
}

func TestListExpensesMissingUserIDReturnsUnauthorized(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := getWithHeaders("/api/v1/expenses", nil)
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusUnauthorized, "Unauthorized")
}

func TestGetOneExpenseSuccessReturnsOK(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	expense := seedExpense(t, 1, "Lunch")

	rec := getWithHeaders("/api/v1/expenses/1", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !response.Success {
		t.Fatal("expected success to be true")
	}
	if response.Message != "Expense retrieved" {
		t.Fatalf("expected message %q, got %q", "Expense retrieved", response.Message)
	}
	if response.Data.ID != expense.ID || response.Data.Title != "Lunch" {
		t.Fatalf("unexpected expense response data: %+v", response.Data)
	}
}

func TestGetOneExpenseReturnsOnlyOwnersExpense(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedAuthenticatedUser(t, 2)
	seedExpense(t, 1, "Lunch")
	expense := seedExpense(t, 2, "Bus")

	rec := getWithHeaders("/api/v1/expenses/2", map[string]string{"X-User-ID": "2"})
	response := decodeExpenseResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if response.Data.ID != expense.ID || response.Data.Title != "Bus" {
		t.Fatalf("expected owner's expense, got %+v", response.Data)
	}
}

func TestGetOneAnotherUsersExpenseReturnsNotFound(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedAuthenticatedUser(t, 2)
	seedExpense(t, 2, "Bus")

	rec := getWithHeaders("/api/v1/expenses/1", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusNotFound, "Expense not found")
}

func TestGetOneMissingExpenseReturnsNotFound(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := getWithHeaders("/api/v1/expenses/99", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusNotFound, "Expense not found")
}

func TestGetOneInvalidIDReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := getWithHeaders("/api/v1/expenses/abc", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid expense ID")
}

func TestGetOneMissingUserIDReturnsUnauthorized(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := getWithHeaders("/api/v1/expenses/1", nil)
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusUnauthorized, "Unauthorized")
}

func TestUpdateExpenseSuccessReturnsOK(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := putJSONWithHeaders("/api/v1/expenses/1", updatedExpenseRequestBody(), map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !response.Success {
		t.Fatal("expected success to be true")
	}
	if response.Message != "Expense updated successfully" {
		t.Fatalf("expected message %q, got %q", "Expense updated successfully", response.Message)
	}
	if response.Data.ID != 1 || response.Data.Title != "Dinner" || response.Data.Amount != 500 || response.Data.Category != "Food" || response.Data.Note != "Family dinner" || response.Data.ExpenseDate != "2025-06-11" {
		t.Fatalf("unexpected expense response data: %+v", response.Data)
	}
}

func TestUpdateExpensePreservesAuthenticatedOwnerUserID(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := putJSONWithHeaders("/api/v1/expenses/1", updatedExpenseRequestBody(), map[string]string{"X-User-ID": "1"})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expense, err := models.GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}
	if expense.UserID != 1 {
		t.Fatalf("expected user ID 1 to be preserved, got %d", expense.UserID)
	}
}

func TestUpdateExpensePreservesCreatedAt(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	originalExpense := seedExpense(t, 1, "Lunch")

	rec := putJSONWithHeaders("/api/v1/expenses/1", updatedExpenseRequestBody(), map[string]string{"X-User-ID": "1"})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	updatedExpense, err := models.GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}
	if updatedExpense.CreatedAt != originalExpense.CreatedAt {
		t.Fatalf("expected CreatedAt %q, got %q", originalExpense.CreatedAt, updatedExpense.CreatedAt)
	}
}

func TestUpdateAnotherUsersExpenseReturnsNotFound(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedAuthenticatedUser(t, 2)
	seedExpense(t, 2, "Bus")

	rec := putJSONWithHeaders("/api/v1/expenses/1", updatedExpenseRequestBody(), map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusNotFound, "Expense not found")
}

func TestUpdateMissingExpenseReturnsNotFound(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := putJSONWithHeaders("/api/v1/expenses/99", updatedExpenseRequestBody(), map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusNotFound, "Expense not found")
}

func TestUpdateInvalidIDReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := putJSONWithHeaders("/api/v1/expenses/abc", updatedExpenseRequestBody(), map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid expense ID")
}

func TestUpdateMissingUserIDReturnsUnauthorized(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := putJSONWithHeaders("/api/v1/expenses/1", updatedExpenseRequestBody(), nil)
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusUnauthorized, "Unauthorized")
}

func TestUpdateInvalidJSONReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := putJSONWithHeaders("/api/v1/expenses/1", `{"title":"Dinner"`, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid request body")
}

func TestUpdateMissingTitleReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	body := `{"amount":500.00,"category":"Food","note":"Family dinner","expense_date":"2025-06-11"}`
	rec := putJSONWithHeaders("/api/v1/expenses/1", body, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Title is required")
}

func TestUpdateZeroAmountReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	body := `{"title":"Dinner","amount":0,"category":"Food","note":"Family dinner","expense_date":"2025-06-11"}`
	rec := putJSONWithHeaders("/api/v1/expenses/1", body, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Amount must be positive")
}

func TestUpdateInvalidCategoryReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	body := `{"title":"Dinner","amount":500.00,"category":"Travel","note":"Family dinner","expense_date":"2025-06-11"}`
	rec := putJSONWithHeaders("/api/v1/expenses/1", body, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid category")
}

func TestUpdateInvalidExpenseDateReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	body := `{"title":"Dinner","amount":500.00,"category":"Food","note":"Family dinner","expense_date":"2025/06/11"}`
	rec := putJSONWithHeaders("/api/v1/expenses/1", body, map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid expense date format")
}

func TestDeleteExpenseSuccessReturnsOK(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := deleteWithHeaders("/api/v1/expenses/1", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !response.Success {
		t.Fatal("expected success to be true")
	}
	if response.Message != "Expense deleted successfully" {
		t.Fatalf("expected message %q, got %q", "Expense deleted successfully", response.Message)
	}
}

func TestDeleteExpenseRemovesExpenseFromCSV(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := deleteWithHeaders("/api/v1/expenses/1", map[string]string{"X-User-ID": "1"})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expense, err := models.GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}
	if expense != nil {
		t.Fatalf("expected expense to be deleted, got %+v", expense)
	}
}

func TestDeleteAnotherUsersExpenseReturnsNotFound(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedAuthenticatedUser(t, 2)
	seedExpense(t, 2, "Bus")

	rec := deleteWithHeaders("/api/v1/expenses/1", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusNotFound, "Expense not found")
}

func TestDeleteMissingExpenseReturnsNotFound(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := deleteWithHeaders("/api/v1/expenses/99", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusNotFound, "Expense not found")
}

func TestDeleteInvalidIDReturnsBadRequest(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)

	rec := deleteWithHeaders("/api/v1/expenses/abc", map[string]string{"X-User-ID": "1"})
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusBadRequest, "Invalid expense ID")
}

func TestDeleteMissingUserIDReturnsUnauthorized(t *testing.T) {
	useTempUserAndExpenseCSVConfig(t)
	seedAuthenticatedUser(t, 1)
	seedExpense(t, 1, "Lunch")

	rec := deleteWithHeaders("/api/v1/expenses/1", nil)
	response := decodeExpenseResponse(t, rec)

	assertExpenseErrorResponse(t, rec.Code, response, http.StatusUnauthorized, "Unauthorized")
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

func seedExpense(t *testing.T, userID int, title string) models.Expense {
	t.Helper()

	expense := &models.Expense{
		UserID:      userID,
		Title:       title,
		Amount:      350.50,
		Category:    "Food",
		Note:        "Team lunch",
		ExpenseDate: "2025-06-10",
	}
	if err := models.CreateExpense(expense); err != nil {
		t.Fatalf("expected expense to be seeded: %v", err)
	}

	return *expense
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

func putJSONWithHeaders(path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(rec, req)

	return rec
}

func getWithHeaders(path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()

	beego.BeeApp.Handlers.ServeHTTP(rec, req)

	return rec
}

func deleteWithHeaders(path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodDelete, path, nil)
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

func decodeExpenseListResponse(t *testing.T, rec *httptest.ResponseRecorder) expenseListResponse {
	t.Helper()

	var response expenseListResponse
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

func updatedExpenseRequestBody() string {
	return `{"title":"Dinner","amount":500.00,"category":"Food","note":"Family dinner","expense_date":"2025-06-11"}`
}
