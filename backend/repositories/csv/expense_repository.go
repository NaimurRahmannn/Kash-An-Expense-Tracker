package csv

import "backend/models"

// ExpenseRepository adapts the existing CSV-backed expense model functions.
type ExpenseRepository struct{}

// NewExpenseRepository creates a CSV expense repository adapter.
func NewExpenseRepository() *ExpenseRepository {
	return &ExpenseRepository{}
}

// GetExpensesByUserID returns all expenses for the user from CSV storage.
func (r *ExpenseRepository) GetExpensesByUserID(userID int) ([]models.Expense, error) {
	return models.GetExpensesByUserID(userID)
}

// GetExpenseByID returns an expense by ID and owner from CSV storage.
func (r *ExpenseRepository) GetExpenseByID(id int, userID int) (*models.Expense, error) {
	return models.GetExpenseByID(id, userID)
}

// CreateExpense stores an expense in CSV storage.
func (r *ExpenseRepository) CreateExpense(expense *models.Expense) error {
	return models.CreateExpense(expense)
}

// UpdateExpense updates an expense in CSV storage.
func (r *ExpenseRepository) UpdateExpense(expense *models.Expense) error {
	return models.UpdateExpense(expense)
}

// DeleteExpense removes an expense from CSV storage.
func (r *ExpenseRepository) DeleteExpense(id int, userID int) error {
	return models.DeleteExpense(id, userID)
}

// GetNextExpenseID returns the next available CSV expense ID.
func (r *ExpenseRepository) GetNextExpenseID() int {
	return models.GetNextExpenseID()
}
