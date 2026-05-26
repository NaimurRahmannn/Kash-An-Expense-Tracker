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

// ExpenseQueryParams represents supported query parameters for listing expenses.
type ExpenseQueryParams struct {
	Category  string
	DateFrom  string
	DateTo    string
	SortBy    string
	SortOrder string
	Page      int
	Limit     int
}

// SummaryQueryParams represents required query parameters for expense summaries.
type SummaryQueryParams struct {
	DateFrom string
	DateTo   string
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

// ValidateExpenseQueryParams validates expense list query parameters.
func ValidateExpenseQueryParams(params ExpenseQueryParams) string {
	if strings.TrimSpace(params.Category) != "" && !IsAllowedCategory(params.Category) {
		return "Invalid category"
	}
	if strings.TrimSpace(params.DateFrom) != "" && !IsValidExpenseDate(params.DateFrom) {
		return "Invalid date_from format"
	}
	if strings.TrimSpace(params.DateTo) != "" && !IsValidExpenseDate(params.DateTo) {
		return "Invalid date_to format"
	}
	if isDateFromAfterDateTo(params.DateFrom, params.DateTo) {
		return "date_from cannot be after date_to"
	}
	if strings.TrimSpace(params.SortBy) != "" && !IsValidSortBy(params.SortBy) {
		return "Invalid sort_by parameter"
	}
	if strings.TrimSpace(params.SortOrder) != "" && !IsValidSortOrder(params.SortOrder) {
		return "Invalid sort_order parameter"
	}
	if params.Page <= 0 {
		return "Invalid page parameter"
	}
	if params.Limit <= 0 {
		return "Invalid limit parameter"
	}

	return ""
}

// ValidateSummaryQueryParams validates expense summary query parameters.
func ValidateSummaryQueryParams(params SummaryQueryParams) string {
	if strings.TrimSpace(params.DateFrom) == "" {
		return "date_from is required"
	}
	if strings.TrimSpace(params.DateTo) == "" {
		return "date_to is required"
	}
	if !IsValidExpenseDate(params.DateFrom) {
		return "Invalid date_from format"
	}
	if !IsValidExpenseDate(params.DateTo) {
		return "Invalid date_to format"
	}
	if isDateFromAfterDateTo(params.DateFrom, params.DateTo) {
		return "date_from cannot be after date_to"
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

// IsValidSortBy reports whether the sort field is supported.
func IsValidSortBy(sortBy string) bool {
	switch strings.TrimSpace(sortBy) {
	case "amount", "expense_date":
		return true
	default:
		return false
	}
}

// IsValidSortOrder reports whether the sort order is supported.
func IsValidSortOrder(sortOrder string) bool {
	switch strings.TrimSpace(sortOrder) {
	case "asc", "desc":
		return true
	default:
		return false
	}
}

// IsValidExpenseDate reports whether the date is valid YYYY-MM-DD.
func IsValidExpenseDate(date string) bool {
	_, err := time.Parse("2006-01-02", strings.TrimSpace(date))
	return err == nil
}

func isDateFromAfterDateTo(dateFrom string, dateTo string) bool {
	normalizedDateFrom := strings.TrimSpace(dateFrom)
	normalizedDateTo := strings.TrimSpace(dateTo)
	if normalizedDateFrom == "" || normalizedDateTo == "" {
		return false
	}
	if !IsValidExpenseDate(normalizedDateFrom) || !IsValidExpenseDate(normalizedDateTo) {
		return false
	}

	return normalizedDateFrom > normalizedDateTo
}
