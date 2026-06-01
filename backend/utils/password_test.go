package utils

import (
	"strings"
	"testing"
)

func TestHashPasswordReturnsNonEmptyHash(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("expected hash to succeed: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "secret123" {
		t.Fatal("expected hash to differ from plain password")
	}
	if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
		t.Fatalf("expected bcrypt hash prefix, got %q", hash)
	}
}

func TestHashPasswordReturnsErrorForEmptyPassword(t *testing.T) {
	if _, err := HashPassword(""); err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestCheckPasswordHashValidatesCorrectPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("expected hash to succeed: %v", err)
	}

	if !CheckPasswordHash("secret123", hash) {
		t.Fatal("expected password to be valid")
	}
}

func TestCheckPasswordHashRejectsWrongPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("expected hash to succeed: %v", err)
	}

	if CheckPasswordHash("wrongpassword", hash) {
		t.Fatal("expected password to be invalid")
	}
}

func TestCheckPasswordHashRejectsEmptyPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("expected hash to succeed: %v", err)
	}

	if CheckPasswordHash("", hash) {
		t.Fatal("expected empty password to be invalid")
	}
}

func TestCheckPasswordHashRejectsEmptyHash(t *testing.T) {
	if CheckPasswordHash("secret123", "") {
		t.Fatal("expected empty hash to be invalid")
	}
}
