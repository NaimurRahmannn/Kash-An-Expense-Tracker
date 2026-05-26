package models

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"backend/utils"

	beego "github.com/beego/beego/v2/server/web"
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

func TestGetExpensesByUserIDReturnsErrorForInvalidAmount(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "1", "Lunch", "invalid", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	})

	if _, err := GetExpensesByUserID(1); err == nil {
		t.Fatal("expected invalid amount to return an error")
	}
}

func TestGetExpensesByUserIDReturnsErrorsForMalformedRows(t *testing.T) {
	tests := []struct {
		name string
		row  []string
	}{
		{
			name: "invalid id",
			row:  []string{"bad", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
		},
		{
			name: "invalid user id",
			row:  []string{"1", "bad", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
		},
		{
			name: "invalid amount",
			row:  []string{"1", "1", "Lunch", "bad", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
		},
		{
			name: "too few columns",
			row:  []string{"1", "1", "Lunch"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := useTempExpenseCSVPath(t)
			writeExpenseRows(t, filePath, [][]string{
				expenseCSVHeader,
				tt.row,
			})

			if _, err := GetExpensesByUserID(1); err == nil {
				t.Fatal("expected malformed expense row to return an error")
			}
		})
	}
}

func TestGetExpensesByUserIDReturnsErrorWhenCSVPathIsDirectory(t *testing.T) {
	previous := expenseCSVFilePathOverride
	expenseCSVFilePathOverride = t.TempDir()
	t.Cleanup(func() {
		expenseCSVFilePathOverride = previous
	})

	if _, err := GetExpensesByUserID(1); err == nil {
		t.Fatal("expected directory expense CSV path to return an error")
	}
}

func TestReadExpenseRowsReturnsEnsureFileError(t *testing.T) {
	parentPath := filepath.Join(t.TempDir(), "parent")
	if err := os.WriteFile(parentPath, []byte("not a directory"), 0644); err != nil {
		t.Fatalf("expected parent file to be written: %v", err)
	}

	if _, err := readExpenseRows(filepath.Join(parentPath, "expenses.csv")); err == nil {
		t.Fatal("expected invalid parent path to return an error")
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

func TestCreateExpenseRequiresExpense(t *testing.T) {
	useTempExpenseCSVPath(t)

	if err := CreateExpense(nil); err == nil {
		t.Fatal("expected nil expense to return an error")
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

func TestGetExpensesByUserIDReturnsEmptySliceForUnknownUser(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	})

	expenses, err := GetExpensesByUserID(99)
	if err != nil {
		t.Fatalf("expected expenses to load: %v", err)
	}
	if len(expenses) != 0 {
		t.Fatalf("expected no expenses for unknown user, got %+v", expenses)
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

func TestGetExpenseByIDReturnsNilForMissingExpense(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	})

	expense, err := GetExpenseByID(99, 1)
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

func TestUpdateExpenseReturnsNotFoundWhenExpenseDoesNotExist(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	})

	err := UpdateExpense(&Expense{
		ID:          99,
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

func TestDeleteExpenseReturnsNotFoundWhenExpenseDoesNotExist(t *testing.T) {
	filePath := useTempExpenseCSVPath(t)
	writeExpenseRows(t, filePath, [][]string{
		expenseCSVHeader,
		{"1", "1", "Lunch", "350.50", "Food", "Team lunch", "2025-06-10", "2025-06-10T14:30:00Z"},
	})

	err := DeleteExpense(99, 1)
	if err == nil {
		t.Fatal("expected expense not found error")
	}
	if err.Error() != "expense not found" {
		t.Fatalf("expected expense not found error, got %q", err.Error())
	}
}

func TestExpenseModelUsesConfiguredCSVPath(t *testing.T) {
	previousOverride := expenseCSVFilePathOverride
	expenseCSVFilePathOverride = ""

	previousPath := beego.AppConfig.DefaultString("csv_expense_file", "data/expenses.csv")
	filePath := filepath.Join(t.TempDir(), "configured-expenses.csv")
	if err := beego.AppConfig.Set("csv_expense_file", filePath); err != nil {
		t.Fatalf("expected csv_expense_file config to be set: %v", err)
	}

	t.Cleanup(func() {
		expenseCSVFilePathOverride = previousOverride
		_ = beego.AppConfig.Set("csv_expense_file", previousPath)
	})

	if err := CreateExpense(&Expense{
		UserID:      1,
		Title:       "Config Lunch",
		Amount:      350.50,
		Category:    "Food",
		Note:        "Team lunch",
		ExpenseDate: "2025-06-10",
	}); err != nil {
		t.Fatalf("expected expense to be created using configured path: %v", err)
	}

	rows, err := utils.ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected configured expenses CSV to be readable: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected header and one expense row, got %d rows", len(rows))
	}
}

func TestGetExpenseCSVFilePathUsesOverride(t *testing.T) {
	previous := expenseCSVFilePathOverride
	filePath := filepath.Join(t.TempDir(), "override-expenses.csv")
	expenseCSVFilePathOverride = filePath
	t.Cleanup(func() {
		expenseCSVFilePathOverride = previous
	})

	got, err := getExpenseCSVFilePath()
	if err != nil {
		t.Fatalf("expected override path to load: %v", err)
	}
	if got != filePath {
		t.Fatalf("expected override path %q, got %q", filePath, got)
	}
}

func TestGetExpenseCSVFilePathReturnsErrorWhenConfigBlank(t *testing.T) {
	previousOverride := expenseCSVFilePathOverride
	expenseCSVFilePathOverride = ""

	previousPath := beego.AppConfig.DefaultString("csv_expense_file", "data/expenses.csv")
	if err := beego.AppConfig.Set("csv_expense_file", "  "); err != nil {
		t.Fatalf("expected csv_expense_file config to be set: %v", err)
	}

	t.Cleanup(func() {
		expenseCSVFilePathOverride = previousOverride
		_ = beego.AppConfig.Set("csv_expense_file", previousPath)
	})

	if _, err := getExpenseCSVFilePath(); err == nil {
		t.Fatal("expected blank csv_expense_file config to return an error")
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
