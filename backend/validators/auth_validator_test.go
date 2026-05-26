package validators

import "testing"

func TestValidateRegisterInput(t *testing.T) {
	tests := []struct {
		name        string
		input       RegisterInput
		wantMessage string
	}{
		{
			name:        "valid register input",
			input:       RegisterInput{Name: "John Doe", Email: "john@example.com", Password: "secret123"},
			wantMessage: "",
		},
		{
			name:        "missing name",
			input:       RegisterInput{Email: "john@example.com", Password: "secret123"},
			wantMessage: "Name is required",
		},
		{
			name:        "missing email",
			input:       RegisterInput{Name: "John Doe", Password: "secret123"},
			wantMessage: "Email is required",
		},
		{
			name:        "invalid email",
			input:       RegisterInput{Name: "John Doe", Email: "invalid-email", Password: "secret123"},
			wantMessage: "Invalid email format",
		},
		{
			name:        "missing password",
			input:       RegisterInput{Name: "John Doe", Email: "john@example.com"},
			wantMessage: "Password is required",
		},
		{
			name:        "short password",
			input:       RegisterInput{Name: "John Doe", Email: "john@example.com", Password: "12345"},
			wantMessage: "Password must be at least 6 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidationMessage(t, ValidateRegisterInput(tt.input), tt.wantMessage)
		})
	}
}

func TestValidateLoginInput(t *testing.T) {
	tests := []struct {
		name        string
		input       LoginInput
		wantMessage string
	}{
		{
			name:        "valid login input",
			input:       LoginInput{Email: "john@example.com", Password: "secret123"},
			wantMessage: "",
		},
		{
			name:        "missing email",
			input:       LoginInput{Password: "secret123"},
			wantMessage: "Email is required",
		},
		{
			name:        "invalid email",
			input:       LoginInput{Email: "invalid-email", Password: "secret123"},
			wantMessage: "Invalid email format",
		},
		{
			name:        "missing password",
			input:       LoginInput{Email: "john@example.com"},
			wantMessage: "Password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertValidationMessage(t, ValidateLoginInput(tt.input), tt.wantMessage)
		})
	}
}

func assertValidationMessage(t *testing.T, actual string, expected string) {
	t.Helper()

	if actual != expected {
		t.Fatalf("expected validation message %q, got %q", expected, actual)
	}
}
