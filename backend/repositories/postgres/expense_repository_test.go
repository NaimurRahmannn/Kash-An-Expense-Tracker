package postgres

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"backend/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewExpenseRepository(t *testing.T) {
	repository, _, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	if repository == nil {
		t.Fatal("expected repository to be non-nil")
	}
}

func TestExpenseRepositoryNilDBBehavior(t *testing.T) {
	repository := NewExpenseRepository(nil)

	tests := []struct {
		name string
		run  func() error
	}{
		{
			name: "GetExpensesByUserID",
			run: func() error {
				_, err := repository.GetExpensesByUserID(1)
				return err
			},
		},
		{
			name: "GetExpenseByID",
			run: func() error {
				_, err := repository.GetExpenseByID(1, 1)
				return err
			},
		},
		{
			name: "CreateExpense",
			run: func() error {
				return repository.CreateExpense(validPostgresExpense())
			},
		},
		{
			name: "UpdateExpense",
			run: func() error {
				return repository.UpdateExpense(validPostgresExpense())
			},
		},
		{
			name: "DeleteExpense",
			run: func() error {
				return repository.DeleteExpense(1, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if err == nil {
				t.Fatal("expected nil DB error")
			}
			if !strings.Contains(err.Error(), "database connection is nil") {
				t.Fatalf("expected nil DB error, got %q", err.Error())
			}
		})
	}

	if nextID := repository.GetNextExpenseID(); nextID != 1 {
		t.Fatalf("expected next ID 1 for nil DB, got %d", nextID)
	}
}

func TestExpenseRepositoryGetExpensesByUserIDReturnsExpenses(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()
	expenseDate := fixedExpenseDate()
	createdAt := fixedPostgresTime()

	rows := sqlmock.NewRows(expenseColumns()).
		AddRow(2, 1, "Dinner", 500.00, "Food", "Family dinner", expenseDate, createdAt).
		AddRow(1, 1, "Lunch", 350.50, "Food", "Team lunch", expenseDate.AddDate(0, 0, -1), createdAt.Add(-time.Hour))
	mock.ExpectQuery(regexp.QuoteMeta(selectExpensesByUserIDQuery)).
		WithArgs(1).
		WillReturnRows(rows)

	expenses, err := repository.GetExpensesByUserID(1)
	if err != nil {
		t.Fatalf("expected expenses to load: %v", err)
	}
	if len(expenses) != 2 {
		t.Fatalf("expected 2 expenses, got %d", len(expenses))
	}
	if expenses[0].ID != 2 || expenses[0].ExpenseDate != "2025-06-10" || expenses[0].CreatedAt != createdAt.UTC().Format(time.RFC3339) {
		t.Fatalf("unexpected first expense: %+v", expenses[0])
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryGetExpensesByUserIDReturnsEmptySlice(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	rows := sqlmock.NewRows(expenseColumns())
	mock.ExpectQuery(regexp.QuoteMeta(selectExpensesByUserIDQuery)).
		WithArgs(1).
		WillReturnRows(rows)

	expenses, err := repository.GetExpensesByUserID(1)
	if err != nil {
		t.Fatalf("expected expenses to load: %v", err)
	}
	if expenses == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(expenses) != 0 {
		t.Fatalf("expected no expenses, got %d", len(expenses))
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryGetExpensesByUserIDReturnsRowsError(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	rows := sqlmock.NewRows(expenseColumns()).
		AddRow(1, 1, "Lunch", 350.50, "Food", "Team lunch", fixedExpenseDate(), fixedPostgresTime()).
		RowError(0, errors.New("row scan failed"))
	mock.ExpectQuery(regexp.QuoteMeta(selectExpensesByUserIDQuery)).
		WithArgs(1).
		WillReturnRows(rows)

	if _, err := repository.GetExpensesByUserID(1); err == nil {
		t.Fatal("expected rows error")
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryGetExpenseByIDReturnsExpense(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	rows := sqlmock.NewRows(expenseColumns()).
		AddRow(1, 1, "Lunch", 350.50, "Food", "Team lunch", fixedExpenseDate(), fixedPostgresTime())
	mock.ExpectQuery(regexp.QuoteMeta(selectExpenseByIDQuery)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	expense, err := repository.GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("expected expense lookup to succeed: %v", err)
	}
	if expense == nil || expense.ID != 1 || expense.UserID != 1 || expense.ExpenseDate != "2025-06-10" {
		t.Fatalf("expected matching expense, got %+v", expense)
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryGetExpenseByIDReturnsNilForMissingExpense(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(selectExpenseByIDQuery)).
		WithArgs(99, 1).
		WillReturnError(sql.ErrNoRows)

	expense, err := repository.GetExpenseByID(99, 1)
	if err != nil {
		t.Fatalf("expected missing expense to return nil error: %v", err)
	}
	if expense != nil {
		t.Fatalf("expected nil expense, got %+v", expense)
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryCreateExpenseAssignsReturnedID(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()
	expense := validPostgresExpense()

	mock.ExpectQuery(regexp.QuoteMeta(insertExpenseQuery)).
		WithArgs(expense.UserID, expense.Title, expense.Amount, expense.Category, expense.Note, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(7))

	if err := repository.CreateExpense(expense); err != nil {
		t.Fatalf("expected expense to be created: %v", err)
	}
	if expense.ID != 7 {
		t.Fatalf("expected returned ID 7, got %d", expense.ID)
	}
	if expense.CreatedAt == "" {
		t.Fatal("expected CreatedAt to be set")
	}
	if _, err := time.Parse(time.RFC3339, expense.CreatedAt); err != nil {
		t.Fatalf("expected RFC3339 CreatedAt, got %q", expense.CreatedAt)
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryCreateExpenseValidationErrors(t *testing.T) {
	repository, _, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	if err := repository.CreateExpense(nil); err == nil || !strings.Contains(err.Error(), "expense is nil") {
		t.Fatalf("expected nil expense error, got %v", err)
	}

	expense := validPostgresExpense()
	expense.ExpenseDate = "invalid-date"
	if err := repository.CreateExpense(expense); err == nil {
		t.Fatal("expected invalid expense date error")
	}

	expense = validPostgresExpense()
	expense.CreatedAt = "invalid-date"
	if err := repository.CreateExpense(expense); err == nil {
		t.Fatal("expected invalid created_at error")
	}
}

func TestExpenseRepositoryCreateExpenseReturnsDatabaseError(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()
	expense := validPostgresExpense()
	expense.CreatedAt = fixedPostgresTime().UTC().Format(time.RFC3339)

	mock.ExpectQuery(regexp.QuoteMeta(insertExpenseQuery)).
		WithArgs(expense.UserID, expense.Title, expense.Amount, expense.Category, expense.Note, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	if err := repository.CreateExpense(expense); err == nil {
		t.Fatal("expected database error")
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryUpdateExpenseSuccess(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()
	expense := validPostgresExpense()
	expense.ID = 1

	mock.ExpectExec(regexp.QuoteMeta(updateExpenseQuery)).
		WithArgs(expense.Title, expense.Amount, expense.Category, expense.Note, sqlmock.AnyArg(), expense.ID, expense.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repository.UpdateExpense(expense); err != nil {
		t.Fatalf("expected expense to update: %v", err)
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryUpdateExpenseReturnsNotFound(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()
	expense := validPostgresExpense()
	expense.ID = 99

	mock.ExpectExec(regexp.QuoteMeta(updateExpenseQuery)).
		WithArgs(expense.Title, expense.Amount, expense.Category, expense.Note, sqlmock.AnyArg(), expense.ID, expense.UserID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repository.UpdateExpense(expense)
	if err == nil {
		t.Fatal("expected expense not found error")
	}
	if !strings.Contains(err.Error(), "expense not found") {
		t.Fatalf("expected expense not found error, got %q", err.Error())
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryUpdateExpenseValidationErrors(t *testing.T) {
	repository, _, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	if err := repository.UpdateExpense(nil); err == nil || !strings.Contains(err.Error(), "expense is nil") {
		t.Fatalf("expected nil expense error, got %v", err)
	}

	expense := validPostgresExpense()
	expense.ExpenseDate = "invalid-date"
	if err := repository.UpdateExpense(expense); err == nil {
		t.Fatal("expected invalid expense date error")
	}
}

func TestExpenseRepositoryUpdateExpenseReturnsDatabaseError(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()
	expense := validPostgresExpense()
	expense.ID = 1

	mock.ExpectExec(regexp.QuoteMeta(updateExpenseQuery)).
		WithArgs(expense.Title, expense.Amount, expense.Category, expense.Note, sqlmock.AnyArg(), expense.ID, expense.UserID).
		WillReturnError(errors.New("update failed"))

	if err := repository.UpdateExpense(expense); err == nil {
		t.Fatal("expected database error")
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryDeleteExpenseSuccess(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(deleteExpenseQuery)).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repository.DeleteExpense(1, 1); err != nil {
		t.Fatalf("expected expense to delete: %v", err)
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryDeleteExpenseReturnsNotFound(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(deleteExpenseQuery)).
		WithArgs(99, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repository.DeleteExpense(99, 1)
	if err == nil {
		t.Fatal("expected expense not found error")
	}
	if !strings.Contains(err.Error(), "expense not found") {
		t.Fatalf("expected expense not found error, got %q", err.Error())
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryDeleteExpenseReturnsDatabaseError(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(deleteExpenseQuery)).
		WithArgs(1, 1).
		WillReturnError(errors.New("delete failed"))

	if err := repository.DeleteExpense(1, 1); err == nil {
		t.Fatal("expected database error")
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryGetNextExpenseIDReturnsNextID(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(selectNextExpenseIDQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"next_id"}).AddRow(12))

	if nextID := repository.GetNextExpenseID(); nextID != 12 {
		t.Fatalf("expected next ID 12, got %d", nextID)
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func TestExpenseRepositoryGetNextExpenseIDReturnsOneOnQueryError(t *testing.T) {
	repository, mock, cleanup := newMockExpenseRepository(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(selectNextExpenseIDQuery)).
		WillReturnError(errors.New("query failed"))

	if nextID := repository.GetNextExpenseID(); nextID != 1 {
		t.Fatalf("expected next ID 1, got %d", nextID)
	}

	assertExpenseRepositoryExpectations(t, mock)
}

func newMockExpenseRepository(t *testing.T) (*ExpenseRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock database: %v", err)
	}

	return NewExpenseRepository(db), mock, func() {
		_ = db.Close()
	}
}

func assertExpenseRepositoryExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func validPostgresExpense() *models.Expense {
	return &models.Expense{
		UserID:      1,
		Title:       "Lunch",
		Amount:      350.50,
		Category:    "Food",
		Note:        "Team lunch",
		ExpenseDate: "2025-06-10",
	}
}

func fixedExpenseDate() time.Time {
	return time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)
}

func expenseColumns() []string {
	return []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}
}
