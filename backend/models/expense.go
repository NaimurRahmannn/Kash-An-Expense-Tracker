package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"backend/config"
	"backend/utils"
)

// Expense represents a user's expense stored in CSV.
type Expense struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
	CreatedAt   string  `json:"created_at"`
}

// AllowedCategories contains the valid categories for expenses.
var AllowedCategories = []string{
	"Food",
	"Transport",
	"Housing",
	"Entertainment",
	"Shopping",
	"Healthcare",
	"Education",
	"Utilities",
	"Other",
}

var (
	expenseCSVHeader           = []string{"id", "user_id", "title", "amount", "category", "note", "expense_date", "created_at"}
	expenseCSVFilePathOverride string
)

// GetExpensesByUserID returns all expenses owned by the provided user.
func GetExpensesByUserID(userID int) ([]Expense, error) {
	expenses, err := getAllExpenses()
	if err != nil {
		return nil, err
	}

	filteredExpenses := make([]Expense, 0)
	for _, expense := range expenses {
		if expense.UserID == userID {
			filteredExpenses = append(filteredExpenses, expense)
		}
	}

	return filteredExpenses, nil
}

// GetExpenseByID returns a matching expense for the provided ID and owner.
func GetExpenseByID(id int, userID int) (*Expense, error) {
	expenses, err := GetExpensesByUserID(userID)
	if err != nil {
		return nil, err
	}

	for index := range expenses {
		if expenses[index].ID == id {
			return &expenses[index], nil
		}
	}

	return nil, nil
}

// CreateExpense appends an expense to the configured CSV file.
func CreateExpense(expense *Expense) error {
	if expense == nil {
		return errors.New("expense is required")
	}

	filePath, err := getExpenseCSVFilePath()
	if err != nil {
		return err
	}

	if err := utils.EnsureFileExists(filePath, expenseCSVHeader); err != nil {
		return err
	}

	if expense.ID == 0 {
		expense.ID = GetNextExpenseID()
	}
	if expense.CreatedAt == "" {
		expense.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	return utils.AppendCSV(filePath, expenseToRow(*expense))
}

// UpdateExpense rewrites the configured CSV with the matching expense updated.
func UpdateExpense(expense *Expense) error {
	if expense == nil {
		return errors.New("expense is required")
	}

	filePath, err := getExpenseCSVFilePath()
	if err != nil {
		return err
	}

	rows, err := readExpenseRows(filePath)
	if err != nil {
		return err
	}

	updatedRows := [][]string{expenseCSVHeader}
	found := false
	for index, row := range rows {
		if index == 0 || isEmptyCSVRow(row) {
			continue
		}

		currentExpense, err := parseExpenseRow(row, index+1)
		if err != nil {
			return err
		}

		if currentExpense.ID == expense.ID && currentExpense.UserID == expense.UserID {
			updatedExpense := *expense
			if updatedExpense.CreatedAt == "" {
				updatedExpense.CreatedAt = currentExpense.CreatedAt
			}
			updatedRows = append(updatedRows, expenseToRow(updatedExpense))
			found = true
			continue
		}

		updatedRows = append(updatedRows, row)
	}

	if !found {
		return errors.New("expense not found")
	}

	return utils.WriteCSV(filePath, updatedRows)
}

// DeleteExpense rewrites the configured CSV without the matching expense.
func DeleteExpense(id int, userID int) error {
	filePath, err := getExpenseCSVFilePath()
	if err != nil {
		return err
	}

	rows, err := readExpenseRows(filePath)
	if err != nil {
		return err
	}

	updatedRows := [][]string{expenseCSVHeader}
	found := false
	for index, row := range rows {
		if index == 0 || isEmptyCSVRow(row) {
			continue
		}

		expense, err := parseExpenseRow(row, index+1)
		if err != nil {
			return err
		}

		if expense.ID == id && expense.UserID == userID {
			found = true
			continue
		}

		updatedRows = append(updatedRows, row)
	}

	if !found {
		return errors.New("expense not found")
	}

	return utils.WriteCSV(filePath, updatedRows)
}

// GetNextExpenseID returns the next globally available expense ID.
func GetNextExpenseID() int {
	expenses, err := getAllExpenses()
	if err != nil {
		return 1
	}

	maxID := 0
	for _, expense := range expenses {
		if expense.ID > maxID {
			maxID = expense.ID
		}
	}

	return maxID + 1
}

func getAllExpenses() ([]Expense, error) {
	filePath, err := getExpenseCSVFilePath()
	if err != nil {
		return nil, err
	}

	rows, err := readExpenseRows(filePath)
	if err != nil {
		return nil, err
	}

	expenses := make([]Expense, 0)
	for index, row := range rows {
		if index == 0 || isEmptyCSVRow(row) {
			continue
		}

		expense, err := parseExpenseRow(row, index+1)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func readExpenseRows(filePath string) ([][]string, error) {
	if err := utils.EnsureFileExists(filePath, expenseCSVHeader); err != nil {
		return nil, err
	}

	return utils.ReadCSV(filePath)
}

func getExpenseCSVFilePath() (string, error) {
	if expenseCSVFilePathOverride != "" {
		return expenseCSVFilePathOverride, nil
	}

	filePath := strings.TrimSpace(config.GetCSVExpenseFile())
	if filePath == "" {
		return "", errors.New("csv_expense_file is not configured")
	}

	return filePath, nil
}

func parseExpenseRow(row []string, rowNumber int) (Expense, error) {
	if len(row) < len(expenseCSVHeader) {
		return Expense{}, fmt.Errorf("invalid expense row %d: expected %d columns, got %d", rowNumber, len(expenseCSVHeader), len(row))
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return Expense{}, fmt.Errorf("invalid expense id at row %d: %w", rowNumber, err)
	}

	userID, err := strconv.Atoi(strings.TrimSpace(row[1]))
	if err != nil {
		return Expense{}, fmt.Errorf("invalid expense user_id at row %d: %w", rowNumber, err)
	}

	amount, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
	if err != nil {
		return Expense{}, fmt.Errorf("invalid expense amount at row %d: %w", rowNumber, err)
	}

	return Expense{
		ID:          id,
		UserID:      userID,
		Title:       row[2],
		Amount:      amount,
		Category:    row[4],
		Note:        row[5],
		ExpenseDate: row[6],
		CreatedAt:   row[7],
	}, nil
}

func expenseToRow(expense Expense) []string {
	return []string{
		strconv.Itoa(expense.ID),
		strconv.Itoa(expense.UserID),
		expense.Title,
		strconv.FormatFloat(expense.Amount, 'f', 2, 64),
		expense.Category,
		expense.Note,
		expense.ExpenseDate,
		expense.CreatedAt,
	}
}
