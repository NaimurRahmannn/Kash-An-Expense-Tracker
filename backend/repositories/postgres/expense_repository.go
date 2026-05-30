package postgres

import (
	"database/sql"
	"errors"
	"time"

	"backend/models"
)

const (
	selectExpensesByUserIDQuery = "SELECT id, user_id, title, amount, category, note, expense_date, created_at FROM expenses WHERE user_id = $1 ORDER BY expense_date DESC, id DESC"
	selectExpenseByIDQuery      = "SELECT id, user_id, title, amount, category, note, expense_date, created_at FROM expenses WHERE id = $1 AND user_id = $2 LIMIT 1"
	insertExpenseQuery          = "INSERT INTO expenses (user_id, title, amount, category, note, expense_date, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	updateExpenseQuery          = "UPDATE expenses SET title = $1, amount = $2, category = $3, note = $4, expense_date = $5 WHERE id = $6 AND user_id = $7"
	deleteExpenseQuery          = "DELETE FROM expenses WHERE id = $1 AND user_id = $2"
	selectNextExpenseIDQuery    = "SELECT COALESCE(MAX(id), 0) + 1 FROM expenses"
)

var (
	// ErrNilExpense is returned when an expense repository receives a nil expense.
	ErrNilExpense = errors.New("expense is nil")
	// ErrExpenseNotFound is returned when an expense does not exist for the owner.
	ErrExpenseNotFound = errors.New("expense not found")
)

// ExpenseRepository stores expenses in Postgres.
type ExpenseRepository struct {
	db *sql.DB
}

// NewExpenseRepository creates a Postgres expense repository.
func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// GetExpensesByUserID returns all expenses for a user ordered by newest expense date.
func (r *ExpenseRepository) GetExpensesByUserID(userID int) ([]models.Expense, error) {
	if r.db == nil {
		return nil, ErrNilDB
	}

	rows, err := r.db.Query(selectExpensesByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]models.Expense, 0)
	for rows.Next() {
		expense, err := scanExpense(rows)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, *expense)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetExpenseByID returns an expense by ID and owner.
func (r *ExpenseRepository) GetExpenseByID(id int, userID int) (*models.Expense, error) {
	if r.db == nil {
		return nil, ErrNilDB
	}

	expense, err := scanExpense(r.db.QueryRow(selectExpenseByIDQuery, id, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return expense, err
}

// CreateExpense inserts an expense and assigns the database-generated ID.
func (r *ExpenseRepository) CreateExpense(expense *models.Expense) error {
	if r.db == nil {
		return ErrNilDB
	}
	if expense == nil {
		return ErrNilExpense
	}
	if expense.CreatedAt == "" {
		expense.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	expenseDate, err := parseExpenseDate(expense.ExpenseDate)
	if err != nil {
		return err
	}
	createdAt, err := time.Parse(time.RFC3339, expense.CreatedAt)
	if err != nil {
		return err
	}

	return r.db.QueryRow(
		insertExpenseQuery,
		expense.UserID,
		expense.Title,
		expense.Amount,
		expense.Category,
		expense.Note,
		expenseDate,
		createdAt,
	).Scan(&expense.ID)
}

// UpdateExpense updates supported expense fields for the matching owner.
func (r *ExpenseRepository) UpdateExpense(expense *models.Expense) error {
	if r.db == nil {
		return ErrNilDB
	}
	if expense == nil {
		return ErrNilExpense
	}

	expenseDate, err := parseExpenseDate(expense.ExpenseDate)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(
		updateExpenseQuery,
		expense.Title,
		expense.Amount,
		expense.Category,
		expense.Note,
		expenseDate,
		expense.ID,
		expense.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrExpenseNotFound
	}

	return nil
}

// DeleteExpense deletes an expense for the matching owner.
func (r *ExpenseRepository) DeleteExpense(id int, userID int) error {
	if r.db == nil {
		return ErrNilDB
	}

	result, err := r.db.Exec(deleteExpenseQuery, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrExpenseNotFound
	}

	return nil
}

// GetNextExpenseID returns the next available expense ID for interface compatibility.
func (r *ExpenseRepository) GetNextExpenseID() int {
	if r.db == nil {
		return 1
	}

	nextID := 1
	if err := r.db.QueryRow(selectNextExpenseIDQuery).Scan(&nextID); err != nil {
		return 1
	}

	return nextID
}

type expenseScanner interface {
	Scan(dest ...any) error
}

func scanExpense(scanner expenseScanner) (*models.Expense, error) {
	var expense models.Expense
	var expenseDate time.Time
	var createdAt time.Time

	if err := scanner.Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Title,
		&expense.Amount,
		&expense.Category,
		&expense.Note,
		&expenseDate,
		&createdAt,
	); err != nil {
		return nil, err
	}

	expense.ExpenseDate = formatExpenseDate(expenseDate)
	expense.CreatedAt = formatPostgresTime(createdAt)
	return &expense, nil
}

func parseExpenseDate(dateString string) (time.Time, error) {
	return time.Parse("2006-01-02", dateString)
}

func formatExpenseDate(value time.Time) string {
	return value.Format("2006-01-02")
}
