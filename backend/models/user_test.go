package models

import (
	"path/filepath"
	"testing"
	"time"

	"backend/utils"
)

func TestGetAllUsersReturnsEmptySliceWhenOnlyHeaderExists(t *testing.T) {
	filePath := useTempUserCSVPath(t)
	ensureUserCSV(t, filePath)

	users, err := GetAllUsers()
	if err != nil {
		t.Fatalf("expected users to load: %v", err)
	}

	if len(users) != 0 {
		t.Fatalf("expected no users, got %d", len(users))
	}
}

func TestCreateUserWritesUserToCSV(t *testing.T) {
	filePath := useTempUserCSVPath(t)
	user := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	if err := CreateUser(user); err != nil {
		t.Fatalf("expected user to be created: %v", err)
	}

	if user.ID != 1 {
		t.Fatalf("expected ID 1, got %d", user.ID)
	}

	if _, err := time.Parse(time.RFC3339, user.CreatedAt); err != nil {
		t.Fatalf("expected CreatedAt to use RFC3339 format: %v", err)
	}

	rows, err := utils.ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected users CSV to be readable: %v", err)
	}

	if len(rows) != 2 {
		t.Fatalf("expected header and one user row, got %d rows", len(rows))
	}

	row := rows[1]
	if row[0] != "1" || row[1] != "John Doe" || row[2] != "john@example.com" || row[3] != "secret123" || row[4] != user.CreatedAt {
		t.Fatalf("unexpected user row: %v", row)
	}
}

func TestGetUserByEmailFindsExistingUser(t *testing.T) {
	useTempUserCSVPath(t)

	if err := CreateUser(&User{Name: "John Doe", Email: "John@Example.com", Password: "secret123"}); err != nil {
		t.Fatalf("expected user to be created: %v", err)
	}

	user, err := GetUserByEmail("john@example.com")
	if err != nil {
		t.Fatalf("expected user lookup to succeed: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be found")
	}

	if user.Email != "John@Example.com" {
		t.Fatalf("expected matching email, got %q", user.Email)
	}
}

func TestGetUserByEmailReturnsNilForMissingUser(t *testing.T) {
	filePath := useTempUserCSVPath(t)
	ensureUserCSV(t, filePath)

	user, err := GetUserByEmail("missing@example.com")
	if err != nil {
		t.Fatalf("expected user lookup to succeed: %v", err)
	}

	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}

func TestGetUserByIDFindsExistingUser(t *testing.T) {
	filePath := useTempUserCSVPath(t)
	ensureUserCSV(t, filePath)

	rows := [][]string{
		userCSVHeader,
		{"1", "John Doe", "john@example.com", "secret123", "2025-06-01T10:30:00Z"},
	}
	if err := utils.WriteCSV(filePath, rows); err != nil {
		t.Fatalf("expected users CSV to be written: %v", err)
	}

	user, err := GetUserByID(1)
	if err != nil {
		t.Fatalf("expected user lookup to succeed: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be found")
	}

	if user.ID != 1 || user.Email != "john@example.com" {
		t.Fatalf("expected matching user, got %+v", user)
	}
}

func TestGetUserByIDReturnsNilForMissingUser(t *testing.T) {
	filePath := useTempUserCSVPath(t)
	ensureUserCSV(t, filePath)

	user, err := GetUserByID(99)
	if err != nil {
		t.Fatalf("expected user lookup to succeed: %v", err)
	}

	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}

func TestGetNextIDReturnsOneForEmptyFile(t *testing.T) {
	filePath := useTempUserCSVPath(t)
	ensureUserCSV(t, filePath)

	nextID := GetNextID()
	if nextID != 1 {
		t.Fatalf("expected next ID 1, got %d", nextID)
	}
}

func TestGetNextIDReturnsMaxIDPlusOne(t *testing.T) {
	filePath := useTempUserCSVPath(t)
	ensureUserCSV(t, filePath)

	rows := [][]string{
		userCSVHeader,
		{"2", "John Doe", "john@example.com", "secret123", "2025-06-01T10:30:00Z"},
		{"5", "Jane Doe", "jane@example.com", "secret456", "2025-06-02T10:30:00Z"},
	}

	if err := utils.WriteCSV(filePath, rows); err != nil {
		t.Fatalf("expected users CSV to be written: %v", err)
	}

	nextID := GetNextID()
	if nextID != 6 {
		t.Fatalf("expected next ID 6, got %d", nextID)
	}
}

func useTempUserCSVPath(t *testing.T) string {
	t.Helper()

	previous := userCSVFilePathOverride
	filePath := filepath.Join(t.TempDir(), "users.csv")
	userCSVFilePathOverride = filePath

	t.Cleanup(func() {
		userCSVFilePathOverride = previous
	})

	return filePath
}

func ensureUserCSV(t *testing.T, filePath string) {
	t.Helper()

	if err := utils.EnsureFileExists(filePath, userCSVHeader); err != nil {
		t.Fatalf("expected users CSV to be created: %v", err)
	}
}
