package csv

import (
	"path/filepath"
	"testing"

	"backend/models"

	beego "github.com/beego/beego/v2/server/web"
)

func TestExpenseRepositoryUsesCSVModelFunctions(t *testing.T) {
	useTempExpenseCSVConfig(t)
	repository := NewExpenseRepository()

	expense := &models.Expense{
		UserID:      1,
		Title:       "Lunch",
		Amount:      350.50,
		Category:    "Food",
		Note:        "Team lunch",
		ExpenseDate: "2025-06-10",
	}
	if err := repository.CreateExpense(expense); err != nil {
		t.Fatalf("expected expense to be created: %v", err)
	}
	if expense.ID != 1 {
		t.Fatalf("expected created expense ID 1, got %d", expense.ID)
	}

	secondExpense := &models.Expense{
		UserID:      2,
		Title:       "Bus",
		Amount:      25,
		Category:    "Transport",
		Note:        "Commute",
		ExpenseDate: "2025-06-11",
	}
	if err := repository.CreateExpense(secondExpense); err != nil {
		t.Fatalf("expected second expense to be created: %v", err)
	}

	expenses, err := repository.GetExpensesByUserID(1)
	if err != nil {
		t.Fatalf("expected expenses by user to load: %v", err)
	}
	if len(expenses) != 1 || expenses[0].ID != expense.ID {
		t.Fatalf("expected one expense for user 1, got %+v", expenses)
	}

	foundExpense, err := repository.GetExpenseByID(expense.ID, 1)
	if err != nil {
		t.Fatalf("expected expense lookup by ID to succeed: %v", err)
	}
	if foundExpense == nil || foundExpense.Title != "Lunch" {
		t.Fatalf("expected created expense by ID, got %+v", foundExpense)
	}

	expense.Title = "Dinner"
	expense.Amount = 500
	expense.Note = "Family dinner"
	expense.ExpenseDate = "2025-06-12"
	if err := repository.UpdateExpense(expense); err != nil {
		t.Fatalf("expected expense to update: %v", err)
	}

	updatedExpense, err := repository.GetExpenseByID(expense.ID, 1)
	if err != nil {
		t.Fatalf("expected updated expense to load: %v", err)
	}
	if updatedExpense.Title != "Dinner" || updatedExpense.Amount != 500 {
		t.Fatalf("expected updated expense, got %+v", updatedExpense)
	}

	if nextID := repository.GetNextExpenseID(); nextID != 3 {
		t.Fatalf("expected next expense ID 3, got %d", nextID)
	}

	if err := repository.DeleteExpense(expense.ID, 1); err != nil {
		t.Fatalf("expected expense to delete: %v", err)
	}

	deletedExpense, err := repository.GetExpenseByID(expense.ID, 1)
	if err != nil {
		t.Fatalf("expected deleted expense lookup to succeed: %v", err)
	}
	if deletedExpense != nil {
		t.Fatalf("expected deleted expense to be nil, got %+v", deletedExpense)
	}
}

func useTempExpenseCSVConfig(t *testing.T) {
	t.Helper()

	previousPath := beego.AppConfig.DefaultString("csv_expense_file", "data/expenses.csv")
	filePath := filepath.Join(t.TempDir(), "expenses.csv")

	if err := beego.AppConfig.Set("csv_expense_file", filePath); err != nil {
		t.Fatalf("expected csv_expense_file test config to be set: %v", err)
	}

	t.Cleanup(func() {
		_ = beego.AppConfig.Set("csv_expense_file", previousPath)
	})
}
