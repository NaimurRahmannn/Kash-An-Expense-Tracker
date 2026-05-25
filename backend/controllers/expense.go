package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"backend/models"
	"backend/validators"

	"github.com/beego/beego/v2/core/logs"
)

// ExpenseController handles expense requests.
type ExpenseController struct {
	BaseController
}

// ExpenseResponse represents expense data returned by the API.
type ExpenseResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
}

// Create creates an expense for the authenticated user.
func (c *ExpenseController) Create() {
	userID, ok := c.getAuthenticatedUserID()
	if !ok {
		c.ErrorResponse(http.StatusUnauthorized, "Unauthorized")
		return
	}

	input, ok := c.parseExpenseInput()
	if !ok {
		return
	}

	expense := &models.Expense{
		UserID:      userID,
		Title:       input.Title,
		Amount:      input.Amount,
		Category:    input.Category,
		Note:        input.Note,
		ExpenseDate: input.ExpenseDate,
	}
	if err := models.CreateExpense(expense); err != nil {
		logs.Error("failed to create expense: %v", err)
		c.ErrorResponse(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.SuccessWithDataAndStatus(http.StatusCreated, "Expense created successfully", toExpenseResponse(*expense))
}

// List returns paginated expenses for the authenticated user.
func (c *ExpenseController) List() {
	userID, ok := c.getAuthenticatedUserID()
	if !ok {
		c.ErrorResponse(http.StatusUnauthorized, "Unauthorized")
		return
	}

	page, ok := c.parsePositiveQueryInt("page", 1, "Invalid page parameter")
	if !ok {
		return
	}
	limit, ok := c.parsePositiveQueryInt("limit", 10, "Invalid limit parameter")
	if !ok {
		return
	}

	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		logs.Error("failed to list expenses for user ID %d: %v", userID, err)
		c.ErrorResponse(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.SuccessWithData("Expenses retrieved", toExpenseResponses(paginateExpenses(expenses, page, limit)))
}

// GetOne returns one expense owned by the authenticated user.
func (c *ExpenseController) GetOne() {
	userID, ok := c.getAuthenticatedUserID()
	if !ok {
		c.ErrorResponse(http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(strings.TrimSpace(c.Ctx.Input.Param(":id")))
	if err != nil || id <= 0 {
		c.ErrorResponse(http.StatusBadRequest, "Invalid expense ID")
		return
	}

	expense, err := models.GetExpenseByID(id, userID)
	if err != nil {
		logs.Error("failed to get expense ID %d for user ID %d: %v", id, userID, err)
		c.ErrorResponse(http.StatusInternalServerError, "Internal server error")
		return
	}
	if expense == nil {
		c.ErrorResponse(http.StatusNotFound, "Expense not found")
		return
	}

	c.SuccessWithData("Expense retrieved", toExpenseResponse(*expense))
}

func (c *ExpenseController) getAuthenticatedUserID() (int, bool) {
	userIDText := strings.TrimSpace(c.Ctx.Input.Header("X-User-ID"))
	if userIDText == "" {
		return 0, false
	}

	userID, err := strconv.Atoi(userIDText)
	if err != nil || userID <= 0 {
		return 0, false
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		logs.Error("failed to authenticate user ID %d: %v", userID, err)
		return 0, false
	}
	if user == nil {
		return 0, false
	}

	return userID, true
}

func (c *ExpenseController) parseExpenseInput() (validators.ExpenseInput, bool) {
	var input validators.ExpenseInput
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		c.ErrorResponse(http.StatusBadRequest, "Invalid request body")
		return input, false
	}

	input = normalizeExpenseInput(input)
	if message := validators.ValidateExpenseInput(input); message != "" {
		c.ErrorResponse(http.StatusBadRequest, message)
		return input, false
	}

	return input, true
}

func (c *ExpenseController) parsePositiveQueryInt(key string, defaultValue int, errorMessage string) (int, bool) {
	valueText := strings.TrimSpace(c.Ctx.Input.Query(key))
	if valueText == "" {
		return defaultValue, true
	}

	value, err := strconv.Atoi(valueText)
	if err != nil || value <= 0 {
		c.ErrorResponse(http.StatusBadRequest, errorMessage)
		return 0, false
	}

	return value, true
}

func normalizeExpenseInput(input validators.ExpenseInput) validators.ExpenseInput {
	return validators.ExpenseInput{
		Title:       strings.TrimSpace(input.Title),
		Amount:      input.Amount,
		Category:    strings.TrimSpace(input.Category),
		Note:        strings.TrimSpace(input.Note),
		ExpenseDate: strings.TrimSpace(input.ExpenseDate),
	}
}

func toExpenseResponse(expense models.Expense) ExpenseResponse {
	return ExpenseResponse{
		ID:          expense.ID,
		Title:       expense.Title,
		Amount:      expense.Amount,
		Category:    expense.Category,
		Note:        expense.Note,
		ExpenseDate: expense.ExpenseDate,
	}
}

func toExpenseResponses(expenses []models.Expense) []ExpenseResponse {
	responses := make([]ExpenseResponse, 0, len(expenses))
	for _, expense := range expenses {
		responses = append(responses, toExpenseResponse(expense))
	}

	return responses
}

func paginateExpenses(expenses []models.Expense, page int, limit int) []models.Expense {
	start := (page - 1) * limit
	if start >= len(expenses) {
		return []models.Expense{}
	}

	end := start + limit
	if end > len(expenses) {
		end = len(expenses)
	}

	return expenses[start:end]
}
