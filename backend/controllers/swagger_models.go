package controllers

// SuccessResponse documents the standard success response envelope.
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse documents the standard error response envelope.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RegisterRequest documents the registration request body.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest documents the login request body.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginData documents the data returned after successful login.
type LoginData struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// LoginSuccessResponse documents the login success response.
type LoginSuccessResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    LoginData `json:"data"`
}

// ExpenseRequest documents create and update expense request bodies.
type ExpenseRequest struct {
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

// ExpenseSuccessResponse documents a single-expense success response.
type ExpenseSuccessResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    ExpenseResponse `json:"data"`
}

// ExpenseListSuccessResponse documents the expense list success response.
type ExpenseListSuccessResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []ExpenseResponse `json:"data"`
}

// ExpenseSummarySuccessResponse documents the expense summary success response.
type ExpenseSummarySuccessResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    ExpenseSummaryResponse `json:"data"`
}
