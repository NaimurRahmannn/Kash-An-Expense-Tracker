package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// BcryptCost controls the hashing cost used for password hashes.
const BcryptCost = bcrypt.DefaultCost

// HashPassword returns a bcrypt hash of the provided password.
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CheckPasswordHash verifies a password against a bcrypt hash.
func CheckPasswordHash(password string, hashedPassword string) bool {
	if password == "" || hashedPassword == "" {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
