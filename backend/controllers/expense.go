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

type createExpenseResponse struct {
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

	c.SuccessWithDataAndStatus(http.StatusCreated, "Expense created successfully", createExpenseResponse{
		ID:          expense.ID,
		Title:       expense.Title,
		Amount:      expense.Amount,
		Category:    expense.Category,
		Note:        expense.Note,
		ExpenseDate: expense.ExpenseDate,
	})
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

func normalizeExpenseInput(input validators.ExpenseInput) validators.ExpenseInput {
	return validators.ExpenseInput{
		Title:       strings.TrimSpace(input.Title),
		Amount:      input.Amount,
		Category:    strings.TrimSpace(input.Category),
		Note:        strings.TrimSpace(input.Note),
		ExpenseDate: strings.TrimSpace(input.ExpenseDate),
	}
}
