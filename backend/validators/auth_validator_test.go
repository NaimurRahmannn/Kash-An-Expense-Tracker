package validators

import "testing"

func TestValidateRegisterInputValid(t *testing.T) {
	input := RegisterInput{Name: "John Doe", Email: "john@example.com", Password: "secret123"}

	if message := ValidateRegisterInput(input); message != "" {
		t.Fatalf("expected valid input, got %q", message)
	}
}

func TestValidateRegisterInputMissingName(t *testing.T) {
	input := RegisterInput{Email: "john@example.com", Password: "secret123"}

	assertValidationMessage(t, ValidateRegisterInput(input), "Name is required")
}

func TestValidateRegisterInputMissingEmail(t *testing.T) {
	input := RegisterInput{Name: "John Doe", Password: "secret123"}

	assertValidationMessage(t, ValidateRegisterInput(input), "Email is required")
}

func TestValidateRegisterInputInvalidEmail(t *testing.T) {
	input := RegisterInput{Name: "John Doe", Email: "invalid-email", Password: "secret123"}

	assertValidationMessage(t, ValidateRegisterInput(input), "Invalid email format")
}

func TestValidateRegisterInputMissingPassword(t *testing.T) {
	input := RegisterInput{Name: "John Doe", Email: "john@example.com"}

	assertValidationMessage(t, ValidateRegisterInput(input), "Password is required")
}

func TestValidateRegisterInputShortPassword(t *testing.T) {
	input := RegisterInput{Name: "John Doe", Email: "john@example.com", Password: "12345"}

	assertValidationMessage(t, ValidateRegisterInput(input), "Password must be at least 6 characters")
}

func TestValidateLoginInputValid(t *testing.T) {
	input := LoginInput{Email: "john@example.com", Password: "secret123"}

	if message := ValidateLoginInput(input); message != "" {
		t.Fatalf("expected valid input, got %q", message)
	}
}

func TestValidateLoginInputMissingEmail(t *testing.T) {
	input := LoginInput{Password: "secret123"}

	assertValidationMessage(t, ValidateLoginInput(input), "Email is required")
}

func TestValidateLoginInputInvalidEmail(t *testing.T) {
	input := LoginInput{Email: "invalid-email", Password: "secret123"}

	assertValidationMessage(t, ValidateLoginInput(input), "Invalid email format")
}

func TestValidateLoginInputMissingPassword(t *testing.T) {
	input := LoginInput{Email: "john@example.com"}

	assertValidationMessage(t, ValidateLoginInput(input), "Password is required")
}

func assertValidationMessage(t *testing.T, actual string, expected string) {
	t.Helper()

	if actual != expected {
		t.Fatalf("expected validation message %q, got %q", expected, actual)
	}
}
