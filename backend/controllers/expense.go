package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"backend/models"
	"backend/repositories"
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

// CategorySummaryResponse represents summary totals for one category.
type CategorySummaryResponse struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}

// ExpenseSummaryResponse represents spending summary data returned by the API.
type ExpenseSummaryResponse struct {
	DateFrom    string                    `json:"date_from"`
	DateTo      string                    `json:"date_to"`
	TotalAmount float64                   `json:"total_amount"`
	TotalCount  int                       `json:"total_count"`
	ByCategory  []CategorySummaryResponse `json:"by_category"`
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
	if err := c.expenseRepository().CreateExpense(expense); err != nil {
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

	params, ok := c.parseExpenseQueryParams()
	if !ok {
		return
	}

	expenses, err := c.expenseRepository().GetExpensesByUserID(userID)
	if err != nil {
		logs.Error("failed to list expenses for user ID %d: %v", userID, err)
		c.ErrorResponse(http.StatusInternalServerError, "Internal server error")
		return
	}

	filteredExpenses := filterExpenses(expenses, params)
	sortedExpenses := sortExpenses(filteredExpenses, params.SortBy, params.SortOrder)
	paginatedExpenses := paginateExpenses(sortedExpenses, params.Page, params.Limit)

	c.SuccessWithData("Expenses retrieved", toExpenseResponses(paginatedExpenses))
}

// Summary returns spending totals for the authenticated user in a date range.
func (c *ExpenseController) Summary() {
	userID, ok := c.getAuthenticatedUserID()
	if !ok {
		c.ErrorResponse(http.StatusUnauthorized, "Unauthorized")
		return
	}

	params, ok := c.parseSummaryQueryParams()
	if !ok {
		return
	}

	expenses, err := c.expenseRepository().GetExpensesByUserID(userID)
	if err != nil {
		logs.Error("failed to summarize expenses for user ID %d: %v", userID, err)
		c.ErrorResponse(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.SuccessWithData("Summary generated", buildExpenseSummary(expenses, params))
}

// GetOne returns one expense owned by the authenticated user.
func (c *ExpenseController) GetOne() {
	userID, ok := c.getAuthenticatedUserID()
	if !ok {
		c.ErrorResponse(http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, ok := c.parseExpenseIDParam()
	if !ok {
		return
	}

	expense, err := c.expenseRepository().GetExpenseByID(id, userID)
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

// Update updates an expense owned by the authenticated user.
func (c *ExpenseController) Update() {
	userID, ok := c.getAuthenticatedUserID()
	if !ok {
		c.ErrorResponse(http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, ok := c.parseExpenseIDParam()
	if !ok {
		return
	}

	existingExpense, ok := c.getExistingExpense(id, userID, "Failed to update expense")
	if !ok {
		return
	}

	input, ok := c.parseExpenseInput()
	if !ok {
		return
	}

	updatedExpense := &models.Expense{
		ID:          existingExpense.ID,
		UserID:      existingExpense.UserID,
		Title:       input.Title,
		Amount:      input.Amount,
		Category:    input.Category,
		Note:        input.Note,
		ExpenseDate: input.ExpenseDate,
		CreatedAt:   existingExpense.CreatedAt,
	}
	if err := c.expenseRepository().UpdateExpense(updatedExpense); err != nil {
		c.handleExpenseWriteError(err, "Failed to update expense")
		return
	}

	c.SuccessWithData("Expense updated successfully", toExpenseResponse(*updatedExpense))
}

// Delete deletes an expense owned by the authenticated user.
func (c *ExpenseController) Delete() {
	userID, ok := c.getAuthenticatedUserID()
	if !ok {
		c.ErrorResponse(http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, ok := c.parseExpenseIDParam()
	if !ok {
		return
	}

	if _, ok := c.getExistingExpense(id, userID, "Failed to delete expense"); !ok {
		return
	}

	if err := c.expenseRepository().DeleteExpense(id, userID); err != nil {
		c.handleExpenseWriteError(err, "Failed to delete expense")
		return
	}

	c.Success("Expense deleted successfully")
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

	user, err := c.userRepository().GetUserByID(userID)
	if err != nil {
		logs.Error("failed to authenticate user ID %d: %v", userID, err)
		return 0, false
	}
	if user == nil {
		return 0, false
	}

	return userID, true
}

func (c *ExpenseController) parseExpenseIDParam() (int, bool) {
	id, err := strconv.Atoi(strings.TrimSpace(c.Ctx.Input.Param(":id")))
	if err != nil || id <= 0 {
		c.ErrorResponse(http.StatusBadRequest, "Invalid expense ID")
		return 0, false
	}

	return id, true
}

func (c *ExpenseController) getExistingExpense(id int, userID int, internalMessage string) (*models.Expense, bool) {
	expense, err := c.expenseRepository().GetExpenseByID(id, userID)
	if err != nil {
		logs.Error("failed to get expense ID %d for user ID %d: %v", id, userID, err)
		c.ErrorResponse(http.StatusInternalServerError, internalMessage)
		return nil, false
	}
	if expense == nil {
		c.ErrorResponse(http.StatusNotFound, "Expense not found")
		return nil, false
	}

	return expense, true
}

func (c *ExpenseController) userRepository() repositories.UserRepository {
	return repositories.NewUserRepository()
}

func (c *ExpenseController) expenseRepository() repositories.ExpenseRepository {
	return repositories.NewExpenseRepository()
}

func (c *ExpenseController) handleExpenseWriteError(err error, internalMessage string) {
	if isExpenseNotFoundError(err) {
		c.ErrorResponse(http.StatusNotFound, "Expense not found")
		return
	}

	logs.Error("%s: %v", strings.ToLower(internalMessage), err)
	c.ErrorResponse(http.StatusInternalServerError, internalMessage)
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

func (c *ExpenseController) parseExpenseQueryParams() (validators.ExpenseQueryParams, bool) {
	params := validators.ExpenseQueryParams{
		Category:  strings.TrimSpace(c.Ctx.Input.Query("category")),
		DateFrom:  strings.TrimSpace(c.Ctx.Input.Query("date_from")),
		DateTo:    strings.TrimSpace(c.Ctx.Input.Query("date_to")),
		SortBy:    strings.TrimSpace(c.Ctx.Input.Query("sort_by")),
		SortOrder: strings.TrimSpace(c.Ctx.Input.Query("sort_order")),
		Page:      parseQueryIntOrDefault(c.Ctx.Input.Query("page"), 1),
		Limit:     parseQueryIntOrDefault(c.Ctx.Input.Query("limit"), 10),
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	if message := validators.ValidateExpenseQueryParams(params); message != "" {
		c.ErrorResponse(http.StatusBadRequest, message)
		return params, false
	}

	return params, true
}

func (c *ExpenseController) parseSummaryQueryParams() (validators.SummaryQueryParams, bool) {
	params := validators.SummaryQueryParams{
		DateFrom: strings.TrimSpace(c.Ctx.Input.Query("date_from")),
		DateTo:   strings.TrimSpace(c.Ctx.Input.Query("date_to")),
	}

	if message := validators.ValidateSummaryQueryParams(params); message != "" {
		c.ErrorResponse(http.StatusBadRequest, message)
		return params, false
	}

	return params, true
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

func filterExpenses(expenses []models.Expense, params validators.ExpenseQueryParams) []models.Expense {
	filteredExpenses := make([]models.Expense, 0, len(expenses))
	for _, expense := range expenses {
		if params.Category != "" && expense.Category != params.Category {
			continue
		}
		if params.DateFrom != "" && expense.ExpenseDate < params.DateFrom {
			continue
		}
		if params.DateTo != "" && expense.ExpenseDate > params.DateTo {
			continue
		}

		filteredExpenses = append(filteredExpenses, expense)
	}

	return filteredExpenses
}

func sortExpenses(expenses []models.Expense, sortBy string, sortOrder string) []models.Expense {
	sortedExpenses := make([]models.Expense, len(expenses))
	copy(sortedExpenses, expenses)

	if sortBy == "" {
		return sortedExpenses
	}

	sort.SliceStable(sortedExpenses, func(i int, j int) bool {
		if sortBy == "amount" {
			return compareExpenseAmount(sortedExpenses[i], sortedExpenses[j], sortOrder)
		}

		return compareExpenseDate(sortedExpenses[i], sortedExpenses[j], sortOrder)
	})

	return sortedExpenses
}

func compareExpenseAmount(left models.Expense, right models.Expense, sortOrder string) bool {
	if sortOrder == "asc" {
		return left.Amount < right.Amount
	}

	return left.Amount > right.Amount
}

func compareExpenseDate(left models.Expense, right models.Expense, sortOrder string) bool {
	if sortOrder == "asc" {
		return left.ExpenseDate < right.ExpenseDate
	}

	return left.ExpenseDate > right.ExpenseDate
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

func buildExpenseSummary(expenses []models.Expense, params validators.SummaryQueryParams) ExpenseSummaryResponse {
	summary := ExpenseSummaryResponse{
		DateFrom:   params.DateFrom,
		DateTo:     params.DateTo,
		ByCategory: []CategorySummaryResponse{},
	}
	categorySummaries := make(map[string]CategorySummaryResponse)

	for _, expense := range expenses {
		if expense.ExpenseDate < params.DateFrom || expense.ExpenseDate > params.DateTo {
			continue
		}

		summary.TotalAmount += expense.Amount
		summary.TotalCount++

		categorySummary := categorySummaries[expense.Category]
		categorySummary.Category = expense.Category
		categorySummary.Total += expense.Amount
		categorySummary.Count++
		categorySummaries[expense.Category] = categorySummary
	}

	categories := make([]string, 0, len(categorySummaries))
	for category := range categorySummaries {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	for _, category := range categories {
		summary.ByCategory = append(summary.ByCategory, categorySummaries[category])
	}

	return summary
}

func isExpenseNotFoundError(err error) bool {
	return err != nil && err.Error() == "expense not found"
}

func parseQueryIntOrDefault(valueText string, defaultValue int) int {
	normalizedValue := strings.TrimSpace(valueText)
	if normalizedValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(normalizedValue)
	if err != nil {
		return 0
	}

	return value
}
