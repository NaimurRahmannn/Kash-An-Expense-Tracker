package validators

import (
	"regexp"
	"strings"
)

// RegisterInput represents the request body for user registration.
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginInput represents the request body for user login.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// ValidateRegisterInput validates registration request data.
func ValidateRegisterInput(input RegisterInput) string {
	if strings.TrimSpace(input.Name) == "" {
		return "Name is required"
	}
	if strings.TrimSpace(input.Email) == "" {
		return "Email is required"
	}
	if !IsValidEmail(input.Email) {
		return "Invalid email format"
	}
	if strings.TrimSpace(input.Password) == "" {
		return "Password is required"
	}
	if len(strings.TrimSpace(input.Password)) < 6 {
		return "Password must be at least 6 characters"
	}

	return ""
}

// ValidateLoginInput validates login request data.
func ValidateLoginInput(input LoginInput) string {
	if strings.TrimSpace(input.Email) == "" {
		return "Email is required"
	}
	if !IsValidEmail(input.Email) {
		return "Invalid email format"
	}
	if strings.TrimSpace(input.Password) == "" {
		return "Password is required"
	}

	return ""
}

// IsValidEmail reports whether an email has a valid format.
func IsValidEmail(email string) bool {
	return emailPattern.MatchString(strings.TrimSpace(email))
}
