package validators

import "testing"

func TestValidateExpenseInputValid(t *testing.T) {
	input := ExpenseInput{
		Title:       "Lunch",
		Amount:      350.50,
		Category:    "Food",
		Note:        "Team lunch",
		ExpenseDate: "2025-06-10",
	}

	if message := ValidateExpenseInput(input); message != "" {
		t.Fatalf("expected valid input, got %q", message)
	}
}

func TestValidateExpenseInputMissingTitle(t *testing.T) {
	input := ExpenseInput{Amount: 350.50, Category: "Food", ExpenseDate: "2025-06-10"}

	assertExpenseValidationMessage(t, ValidateExpenseInput(input), "Title is required")
}

func TestValidateExpenseInputZeroAmount(t *testing.T) {
	input := ExpenseInput{Title: "Lunch", Amount: 0, Category: "Food", ExpenseDate: "2025-06-10"}

	assertExpenseValidationMessage(t, ValidateExpenseInput(input), "Amount must be positive")
}

func TestValidateExpenseInputNegativeAmount(t *testing.T) {
	input := ExpenseInput{Title: "Lunch", Amount: -1, Category: "Food", ExpenseDate: "2025-06-10"}

	assertExpenseValidationMessage(t, ValidateExpenseInput(input), "Amount must be positive")
}

func TestValidateExpenseInputMissingCategory(t *testing.T) {
	input := ExpenseInput{Title: "Lunch", Amount: 350.50, ExpenseDate: "2025-06-10"}

	assertExpenseValidationMessage(t, ValidateExpenseInput(input), "Category is required")
}

func TestValidateExpenseInputInvalidCategory(t *testing.T) {
	input := ExpenseInput{Title: "Lunch", Amount: 350.50, Category: "Travel", ExpenseDate: "2025-06-10"}

	assertExpenseValidationMessage(t, ValidateExpenseInput(input), "Invalid category")
}

func TestValidateExpenseInputMissingExpenseDate(t *testing.T) {
	input := ExpenseInput{Title: "Lunch", Amount: 350.50, Category: "Food"}

	assertExpenseValidationMessage(t, ValidateExpenseInput(input), "Expense date is required")
}

func TestValidateExpenseInputInvalidExpenseDate(t *testing.T) {
	input := ExpenseInput{Title: "Lunch", Amount: 350.50, Category: "Food", ExpenseDate: "2025/06/10"}

	assertExpenseValidationMessage(t, ValidateExpenseInput(input), "Invalid expense date format")
}

func assertExpenseValidationMessage(t *testing.T, actual string, expected string) {
	t.Helper()

	if actual != expected {
		t.Fatalf("expected validation message %q, got %q", expected, actual)
	}
}
