package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/models"
	"backend/repositories"
	"backend/utils"
	"backend/validators"

	"github.com/beego/beego/v2/core/logs"
)

// AuthController handles user registration and login requests.
type AuthController struct {
	BaseController
}

type loginResponse struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// Register creates a new user account.
func (c *AuthController) Register() {
	input, ok := c.parseRegisterInput()
	if !ok {
		return
	}

	userRepo := c.userRepository()
	existingUser, err := userRepo.GetUserByEmail(input.Email)
	if err != nil {
		logs.Error("failed to check existing user: %v", err)
		c.ErrorResponse(http.StatusInternalServerError, "Internal server error")
		return
	}
	if existingUser != nil {
		c.ErrorResponse(http.StatusConflict, "Email already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		logs.Error("failed to hash password: %v", err)
		c.ErrorResponse(http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}
	if err := userRepo.CreateUser(user); err != nil {
		logs.Error("failed to create user: %v", err)
		c.ErrorResponse(http.StatusInternalServerError, "Internal server error")
		return
	}

	c.SuccessWithStatus(http.StatusCreated, "User registered successfully")
}

// Login authenticates a user with email and password.
func (c *AuthController) Login() {
	input, ok := c.parseLoginInput()
	if !ok {
		return
	}

	user, err := c.userRepository().GetUserByEmail(input.Email)
	if err != nil {
		logs.Error("failed to get user by email: %v", err)
		c.ErrorResponse(http.StatusInternalServerError, "Internal server error")
		return
	}
	if user == nil || !utils.CheckPasswordHash(input.Password, user.Password) {
		c.ErrorResponse(http.StatusUnauthorized, "Invalid email or password")
		return
	}

	c.SuccessWithData("Login successful", loginResponse{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	})
}

func (c *AuthController) userRepository() repositories.UserRepository {
	return repositories.NewUserRepository()
}

func (c *AuthController) parseRegisterInput() (validators.RegisterInput, bool) {
	var input validators.RegisterInput
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		c.ErrorResponse(http.StatusBadRequest, "Invalid request body")
		return input, false
	}

	input = normalizeRegisterInput(input)
	if message := validators.ValidateRegisterInput(input); message != "" {
		c.ErrorResponse(http.StatusBadRequest, message)
		return input, false
	}

	return input, true
}

func (c *AuthController) parseLoginInput() (validators.LoginInput, bool) {
	var input validators.LoginInput
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &input); err != nil {
		c.ErrorResponse(http.StatusBadRequest, "Invalid request body")
		return input, false
	}

	input = normalizeLoginInput(input)
	if message := validators.ValidateLoginInput(input); message != "" {
		c.ErrorResponse(http.StatusBadRequest, message)
		return input, false
	}

	return input, true
}

func normalizeRegisterInput(input validators.RegisterInput) validators.RegisterInput {
	return validators.RegisterInput{
		Name:     strings.TrimSpace(input.Name),
		Email:    strings.TrimSpace(input.Email),
		Password: strings.TrimSpace(input.Password),
	}
}

func normalizeLoginInput(input validators.LoginInput) validators.LoginInput {
	return validators.LoginInput{
		Email:    strings.TrimSpace(input.Email),
		Password: strings.TrimSpace(input.Password),
	}
}
