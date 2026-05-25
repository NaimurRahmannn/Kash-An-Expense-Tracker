package models

import (
	"path/filepath"
	"testing"
	"time"

	"backend/utils"
)

func TestGetExpensesByUserIDReturnsEmptySliceWhenOnlyHeaderExists(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	ensureExpenseCSV(t, filePath)

	expenses, err := GetExpensesByUserID(1)
	if err != nil {
		t.Fatalf("expected expenses to load: %v", err)
	}

	if len(expenses) != 0 {
		t.Fatalf("expected no expenses, got %d", len(expenses))
	}
}

func TestCreateExpenseWritesExpenseToCSV(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	expense := &Expense{
		UserID:      1,
		Title:       "Lunch",
		Amount:      350.5,
		Category:    "Food",
		Note:        "Team lunch",
		ExpenseDate: "2025-06-10",
	}

	if err := CreateExpense(expense); err != nil {
		t.Fatalf("expected expense to be created: %v", err)
	}

	if expense.ID != 1 {
		t.Fatalf("expected ID 1, got %d", expense.ID)
	}

	if _, err := time.Parse(time.RFC3339, expense.CreatedAt); err != nil {
		t.Fatalf("expected CreatedAt to use RFC3339 format: %v", err)
	}

	rows, err := utils.ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected expenses CSV to be readable: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("expected header and one expense row, got %d rows", len(rows))
	}

	row := rows[1]
	if row[0] != "1" || row[1] != "1" || row[2] != "Lunch" || row[3] != "350.50" || row[4] != "Food" || row[5] != "Team lunch" || row[6] != "2025-06-10" || row[7] != expense.CreatedAt {
		t.Fatalf("unexpected expense row: %v", row)
	}
}

func TestGetExpensesByUserIDReturnsOnlyRequestedUser(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
		{"2", "2", "Bus", "25.00", "Transport", "Commute", "2025-06-11", "2025-06-11T08:30:00Z"},
	})

	expenses, err := GetExpensesByUserID(1)
	if err != nil {
		t.Fatalf("expected expenses to load: %v", err)
	}

	if len(expenses) != 1 {
		t.Fatalf("expected one expense, got %d", len(expenses))
	}
	if expenses[0].ID != 1 || expenses[0].UserID != 1 {
		t.Fatalf("expected user 1 expense, got %+v", expenses[0])
	}
}

func TestGetExpenseByIDFindsExpenseOnlyForOwner(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	})

	expense, err := GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}

	if expense == nil {
		t.Fatal("expected expense to be found")
	}
	if expense.ID != 1 || expense.UserID != 1 {
		t.Fatalf("expected owner expense, got %+v", expense)
	}
}

func TestGetExpenseByIDReturnsNilForAnotherUsersExpense(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "2", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	})

	expense, err := GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}

	if expense != nil {
		t.Fatalf("expected nil expense, got %+v", expense)
	}
}

func TestGetNextExpenseIDReturnsOneForEmptyFile(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	ensureExpenseCSV(t, filePath)

	nextID := GetNextExpenseID()
	if nextID != 1 {
		t.Fatalf("expected next ID 1, got %d", nextID)
	}
}

func TestGetNextExpenseIDReturnsMaxIDPlusOne(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"2", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
		{"7", "2", "Bus", "25.00", "Transport", "Commute", "2025-06-11", "2025-06-11T08:30:00Z"},
	})

	nextID := GetNextExpenseID()
	if nextID != 8 {
		t.Fatalf("expected next ID 8, got %d", nextID)
	}
}

func TestUpdateExpenseUpdatesMatchingExpense(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	createdAt := "2025-06-10T14:30:00Z"
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", createdAt},
	})

	err := UpdateExpense(&Expense{
		ID:          1,
		UserID:      1,
		Title:       "Dinner",
		Amount:      420,
		Category:    "Food",
		Note:        "Family dinner",
		ExpenseDate: "2025-06-12",
	})
	if err != nil {
		t.Fatalf("expected expense to update: %v", err)
	}

	expense, err := GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}

	if expense.Title != "Dinner" || expense.Amount != 420 || expense.Note != "Family dinner" || expense.ExpenseDate != "2025-06-12" {
		t.Fatalf("expected updated expense, got %+v", expense)
	}
	if expense.CreatedAt != createdAt {
		t.Fatalf("expected CreatedAt to be preserved, got %q", expense.CreatedAt)
	}
}

func TestUpdateExpenseDoesNotUpdateAnotherUsersExpense(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "2", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	})

	err := UpdateExpense(&Expense{
		ID:          1,
		UserID:      1,
		Title:       "Dinner",
		Amount:      420,
		Category:    "Food",
		Note:        "Family dinner",
		ExpenseDate: "2025-06-12",
	})
	if err == nil {
		t.Fatal("expected expense not found error")
	}
	if err.Error() != "expense not found" {
		t.Fatalf("expected expense not found error, got %q", err.Error())
	}

	expense, err := GetExpenseByID(1, 2)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}
	if expense.Title != "Lunch" {
		t.Fatalf("expected another user's expense to remain unchanged, got %+v", expense)
	}
}

func TestDeleteExpenseDeletesMatchingExpense(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
		{"2", "1", "Bus", "25.00", "Transport", "Commute", "2025-06-11", "2025-06-11T08:30:00Z"},
	})

	if err := DeleteExpense(1, 1); err != nil {
		t.Fatalf("expected expense to delete: %v", err)
	}

	expense, err := GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}
	if expense != nil {
		t.Fatalf("expected deleted expense to be nil, got %+v", expense)
	}

	expenses, err := GetExpensesByUserID(1)
	if err != nil {
		t.Fatalf("expected expenses to load: %v", err)
	}
	if len(expenses) != 1 || expenses[0].ID != 2 {
		t.Fatalf("expected remaining expense ID 2, got %+v", expenses)
	}
}

func TestDeleteExpenseDoesNotDeleteAnotherUsersExpense(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "2", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	})

	err := DeleteExpense(1, 1)
	if err == nil {
		t.Fatal("expected expense not found error")
	}
	if err.Error() != "expense not found" {
		t.Fatalf("expected expense not found error, got %q", err.Error())
	}

	expense, err := GetExpenseByID(1, 2)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}
	if expense == nil {
		t.Fatal("expected another user's expense to remain")
	}
}

func useTempExpenseCSVPath(t *testing.T) string {
	t.Helper()

	previous := expenseCSVFilePathOverride
	filePath := filepath.Join(t.TempDir(), "expenses.csv")
	expenseCSVFilePathOverride = filePath

	t.Cleanup(func() {
		expenseCSVFilePathOverride = previous
	})

	return filePath
}

func ensureExpenseCSV(t *testing.T, filePath string) {
	t.Helper()

	if err := utils.EnsureFileExists(filePath, expenseCSVHeader); err != nil {
		t.Fatalf("expected expenses CSV to be created: %v", err)
	}
}

func writeExpenseRows(t *testing.T, filePath string, rows [][]string) {
	t.Helper()

	if err := utils.WriteCSV(filePath, rows); err != nil {
		t.Fatalf("expected expenses CSV to be written: %v", err)
	}
}
