package repositories

import "backend/models"

// ExpenseRepository defines expense persistence behavior across storage drivers.
type ExpenseRepository interface {
	GetExpensesByUserID(userID int) ([]models.Expense, error)
	GetExpenseByID(id int, userID int) (*models.Expense, error)
	CreateExpense(expense *models.Expense) error
	UpdateExpense(expense *models.Expense) error
	DeleteExpense(id int, userID int) error
	GetNextExpenseID() int
}
