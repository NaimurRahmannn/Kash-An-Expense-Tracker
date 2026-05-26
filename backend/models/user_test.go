package models

import (
	"path/filepath"
	"testing"
	"time"

	"backend/utils"

	beego "github.com/beego/beego/v2/server/web"
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

func TestGetAllUsersReturnsErrorForInvalidID(t *testing.T) {
	filePath := useTempUserCSVPath(t)
	rows := [][]string{
		userCSVHeader,
		{"abc", "John Doe", "john@example.com", "secret123", "2025-06-01T10:30:00Z"},
	}
	if err := utils.WriteCSV(filePath, rows); err != nil {
		t.Fatalf("expected users CSV to be written: %v", err)
	}

	if _, err := GetAllUsers(); err == nil {
		t.Fatal("expected invalid user ID to return an error")
	}
}

func TestGetAllUsersReturnsErrorForMalformedRows(t *testing.T) {
	tests := []struct {
		name string
		row  []string
	}{
		{
			name: "invalid id",
			row:  []string{"bad", "John Doe", "john@example.com", "secret123", "2025-06-01T10:30:00Z"},
		},
		{
			name: "too few columns",
			row:  []string{"1", "John Doe"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := useTempUserCSVPath(t)
			if err := utils.WriteCSV(filePath, [][]string{userCSVHeader, tt.row}); err != nil {
				t.Fatalf("expected users CSV to be written: %v", err)
			}

			if _, err := GetAllUsers(); err == nil {
				t.Fatal("expected malformed user row to return an error")
			}
		})
	}
}

func TestGetAllUsersIgnoresEmptyRows(t *testing.T) {
	filePath := useTempUserCSVPath(t)
	rows := [][]string{
		userCSVHeader,
		{"", "", "", "", ""},
		{"   ", "\t", " ", "", ""},
		{"1", "John Doe", "john@example.com", "secret123", "2025-06-01T10:30:00Z"},
	}
	if err := utils.WriteCSV(filePath, rows); err != nil {
		t.Fatalf("expected users CSV to be written: %v", err)
	}

	users, err := GetAllUsers()
	if err != nil {
		t.Fatalf("expected users to load: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected one valid user, got %d users: %+v", len(users), users)
	}
	if users[0].Email != "john@example.com" {
		t.Fatalf("expected valid user to load, got %+v", users[0])
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

func TestCreateUserRequiresUser(t *testing.T) {
	useTempUserCSVPath(t)

	if err := CreateUser(nil); err == nil {
		t.Fatal("expected nil user to return an error")
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

func TestGetUserByEmailFindsEmailCaseInsensitively(t *testing.T) {
	useTempUserCSVPath(t)

	if err := CreateUser(&User{Name: "John Doe", Email: "John@Example.com", Password: "secret123"}); err != nil {
		t.Fatalf("expected user to be created: %v", err)
	}

	tests := []struct {
		name  string
		email string
	}{
		{name: "lowercase email", email: "john@example.com"},
		{name: "uppercase email", email: "JOHN@EXAMPLE.COM"},
		{name: "email with spaces", email: " john@example.com "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := GetUserByEmail(tt.email)
			if err != nil {
				t.Fatalf("expected user lookup to succeed: %v", err)
			}
			if user == nil {
				t.Fatal("expected user to be found")
			}
			if user.Email != "John@Example.com" {
				t.Fatalf("expected matching user email, got %q", user.Email)
			}
		})
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

func TestUserModelUsesConfiguredCSVPath(t *testing.T) {
	previousOverride := userCSVFilePathOverride
	userCSVFilePathOverride = ""

	previousPath := beego.AppConfig.DefaultString("csv_user_file", "data/users.csv")
	filePath := filepath.Join(t.TempDir(), "configured-users.csv")
	if err := beego.AppConfig.Set("csv_user_file", filePath); err != nil {
		t.Fatalf("expected csv_user_file config to be set: %v", err)
	}

	t.Cleanup(func() {
		userCSVFilePathOverride = previousOverride
		_ = beego.AppConfig.Set("csv_user_file", previousPath)
	})

	if err := CreateUser(&User{Name: "Config User", Email: "config@example.com", Password: "secret123"}); err != nil {
		t.Fatalf("expected user to be created using configured path: %v", err)
	}

	rows, err := utils.ReadCSV(filePath)
	if err != nil {
		t.Fatalf("expected configured users CSV to be readable: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected header and one user row, got %d rows", len(rows))
	}
}

func TestGetUserCSVFilePathUsesOverride(t *testing.T) {
	previous := userCSVFilePathOverride
	filePath := filepath.Join(t.TempDir(), "override-users.csv")
	userCSVFilePathOverride = filePath
	t.Cleanup(func() {
		userCSVFilePathOverride = previous
	})

	got, err := getUserCSVFilePath()
	if err != nil {
		t.Fatalf("expected override path to load: %v", err)
	}
	if got != filePath {
		t.Fatalf("expected override path %q, got %q", filePath, got)
	}
}

func TestGetUserCSVFilePathReturnsErrorWhenConfigBlank(t *testing.T) {
	previousOverride := userCSVFilePathOverride
	userCSVFilePathOverride = ""

	previousPath := beego.AppConfig.DefaultString("csv_user_file", "data/users.csv")
	if err := beego.AppConfig.Set("csv_user_file", "  "); err != nil {
		t.Fatalf("expected csv_user_file config to be set: %v", err)
	}

	t.Cleanup(func() {
		userCSVFilePathOverride = previousOverride
		_ = beego.AppConfig.Set("csv_user_file", previousPath)
	})

	if _, err := getUserCSVFilePath(); err == nil {
		t.Fatal("expected blank csv_user_file config to return an error")
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
