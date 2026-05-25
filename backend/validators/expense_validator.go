package validators

import (
	"strings"
	"time"

	"backend/models"
)

// ExpenseInput represents request data for creating or updating an expense.
type ExpenseInput struct {
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

// ValidateExpenseInput validates expense request data.
func ValidateExpenseInput(input ExpenseInput) string {
	if strings.TrimSpace(input.Title) == "" {
		return "Title is required"
	}
	if input.Amount <= 0 {
		return "Amount must be positive"
	}
	if strings.TrimSpace(input.ExpenseDate) == "" {
		return "Expense date is required"
	}
	if !IsValidExpenseDate(input.ExpenseDate) {
		return "Invalid expense date format"
	}
	if strings.TrimSpace(input.Category) == "" {
		return "Category is required"
	}
	if !IsAllowedCategory(input.Category) {
		return "Invalid category"
	}

	return ""
}

// IsAllowedCategory reports whether the category is valid for expenses.
func IsAllowedCategory(category string) bool {
	normalizedCategory := strings.TrimSpace(category)
	for _, allowedCategory := range models.AllowedCategories {
		if normalizedCategory == allowedCategory {
			return true
		}
	}

	return false
}

// IsValidExpenseDate reports whether the date is valid YYYY-MM-DD.
func IsValidExpenseDate(date string) bool {
	_, err := time.Parse("2006-01-02", strings.TrimSpace(date))
	return err == nil
}
