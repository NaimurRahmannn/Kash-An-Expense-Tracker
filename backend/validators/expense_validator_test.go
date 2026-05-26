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

func TestValidateExpenseQueryParamsValid(t *testing.T) {
	params := ExpenseQueryParams{
		Category:  "Food",
		DateFrom:  "2025-06-01",
		DateTo:    "2025-06-30",
		SortBy:    "amount",
		SortOrder: "desc",
		Page:      1,
		Limit:     10,
	}

	if message := ValidateExpenseQueryParams(params); message != "" {
		t.Fatalf("expected valid query params, got %q", message)
	}
}

func TestValidateExpenseQueryParamsInvalidCategory(t *testing.T) {
	params := validExpenseQueryParams()
	params.Category = "Travel"

	assertExpenseValidationMessage(t, ValidateExpenseQueryParams(params), "Invalid category")
}

func TestValidateExpenseQueryParamsValidDateFrom(t *testing.T) {
	params := validExpenseQueryParams()
	params.DateFrom = "2025-06-01"
	params.DateTo = ""

	if message := ValidateExpenseQueryParams(params); message != "" {
		t.Fatalf("expected valid date_from, got %q", message)
	}
}

func TestValidateExpenseQueryParamsInvalidDateFrom(t *testing.T) {
	params := validExpenseQueryParams()
	params.DateFrom = "2025/06/01"

	assertExpenseValidationMessage(t, ValidateExpenseQueryParams(params), "Invalid date_from format")
}

func TestValidateExpenseQueryParamsValidDateTo(t *testing.T) {
	params := validExpenseQueryParams()
	params.DateFrom = ""
	params.DateTo = "2025-06-30"

	if message := ValidateExpenseQueryParams(params); message != "" {
		t.Fatalf("expected valid date_to, got %q", message)
	}
}

func TestValidateExpenseQueryParamsInvalidDateTo(t *testing.T) {
	params := validExpenseQueryParams()
	params.DateTo = "2025/06/30"

	assertExpenseValidationMessage(t, ValidateExpenseQueryParams(params), "Invalid date_to format")
}

func TestValidateExpenseQueryParamsDateFromAfterDateTo(t *testing.T) {
	params := validExpenseQueryParams()
	params.DateFrom = "2025-06-30"
	params.DateTo = "2025-06-01"

	assertExpenseValidationMessage(t, ValidateExpenseQueryParams(params), "date_from cannot be after date_to")
}

func TestValidateExpenseQueryParamsValidSortByAmount(t *testing.T) {
	params := validExpenseQueryParams()
	params.SortBy = "amount"

	if message := ValidateExpenseQueryParams(params); message != "" {
		t.Fatalf("expected valid sort_by amount, got %q", message)
	}
}

func TestValidateExpenseQueryParamsValidSortByExpenseDate(t *testing.T) {
	params := validExpenseQueryParams()
	params.SortBy = "expense_date"

	if message := ValidateExpenseQueryParams(params); message != "" {
		t.Fatalf("expected valid sort_by expense_date, got %q", message)
	}
}

func TestValidateExpenseQueryParamsInvalidSortBy(t *testing.T) {
	params := validExpenseQueryParams()
	params.SortBy = "title"

	assertExpenseValidationMessage(t, ValidateExpenseQueryParams(params), "Invalid sort_by parameter")
}

func TestValidateExpenseQueryParamsValidSortOrderAsc(t *testing.T) {
	params := validExpenseQueryParams()
	params.SortOrder = "asc"

	if message := ValidateExpenseQueryParams(params); message != "" {
		t.Fatalf("expected valid sort_order asc, got %q", message)
	}
}

func TestValidateExpenseQueryParamsValidSortOrderDesc(t *testing.T) {
	params := validExpenseQueryParams()
	params.SortOrder = "desc"

	if message := ValidateExpenseQueryParams(params); message != "" {
		t.Fatalf("expected valid sort_order desc, got %q", message)
	}
}

func TestValidateExpenseQueryParamsInvalidSortOrder(t *testing.T) {
	params := validExpenseQueryParams()
	params.SortOrder = "newest"

	assertExpenseValidationMessage(t, ValidateExpenseQueryParams(params), "Invalid sort_order parameter")
}

func TestValidateExpenseQueryParamsInvalidPage(t *testing.T) {
	params := validExpenseQueryParams()
	params.Page = 0

	assertExpenseValidationMessage(t, ValidateExpenseQueryParams(params), "Invalid page parameter")
}

func TestValidateExpenseQueryParamsInvalidLimit(t *testing.T) {
	params := validExpenseQueryParams()
	params.Limit = 0

	assertExpenseValidationMessage(t, ValidateExpenseQueryParams(params), "Invalid limit parameter")
}

func TestValidateSummaryQueryParamsValid(t *testing.T) {
	params := SummaryQueryParams{DateFrom: "2025-06-01", DateTo: "2025-06-30"}

	if message := ValidateSummaryQueryParams(params); message != "" {
		t.Fatalf("expected valid summary query params, got %q", message)
	}
}

func TestValidateSummaryQueryParamsMissingDateFrom(t *testing.T) {
	params := SummaryQueryParams{DateTo: "2025-06-30"}

	assertExpenseValidationMessage(t, ValidateSummaryQueryParams(params), "date_from is required")
}

func TestValidateSummaryQueryParamsMissingDateTo(t *testing.T) {
	params := SummaryQueryParams{DateFrom: "2025-06-01"}

	assertExpenseValidationMessage(t, ValidateSummaryQueryParams(params), "date_to is required")
}

func TestValidateSummaryQueryParamsInvalidDateFrom(t *testing.T) {
	params := SummaryQueryParams{DateFrom: "2025/06/01", DateTo: "2025-06-30"}

	assertExpenseValidationMessage(t, ValidateSummaryQueryParams(params), "Invalid date_from format")
}

func TestValidateSummaryQueryParamsInvalidDateTo(t *testing.T) {
	params := SummaryQueryParams{DateFrom: "2025-06-01", DateTo: "2025/06/30"}

	assertExpenseValidationMessage(t, ValidateSummaryQueryParams(params), "Invalid date_to format")
}

func TestValidateSummaryQueryParamsDateFromAfterDateTo(t *testing.T) {
	params := SummaryQueryParams{DateFrom: "2025-06-30", DateTo: "2025-06-01"}

	assertExpenseValidationMessage(t, ValidateSummaryQueryParams(params), "date_from cannot be after date_to")
}

func assertExpenseValidationMessage(t *testing.T, actual string, expected string) {
	t.Helper()

	if actual != expected {
		t.Fatalf("expected validation message %q, got %q", expected, actual)
	}
}

func validExpenseQueryParams() ExpenseQueryParams {
	return ExpenseQueryParams{
		Page:      1,
		Limit:     10,
		SortOrder: "desc",
	}
}
