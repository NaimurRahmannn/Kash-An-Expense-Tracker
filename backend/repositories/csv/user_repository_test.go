package csv

import (
	"path/filepath"
	"runtime"
	"testing"

	"backend/models"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	appPath := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	beego.TestBeegoInit(appPath)
}

func TestUserRepositoryUsesCSVModelFunctions(t *testing.T) {
	useTempUserCSVConfig(t)
	repository := NewUserRepository()

	user := &models.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}
	if err := repository.CreateUser(user); err != nil {
		t.Fatalf("expected user to be created: %v", err)
	}

	if user.ID != 1 {
		t.Fatalf("expected created user ID 1, got %d", user.ID)
	}

	userByEmail, err := repository.GetUserByEmail("JOHN@example.com")
	if err != nil {
		t.Fatalf("expected user lookup by email to succeed: %v", err)
	}
	if userByEmail == nil || userByEmail.ID != user.ID {
		t.Fatalf("expected created user by email, got %+v", userByEmail)
	}

	userByID, err := repository.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("expected user lookup by ID to succeed: %v", err)
	}
	if userByID == nil || userByID.Email != user.Email {
		t.Fatalf("expected created user by ID, got %+v", userByID)
	}

	users, err := repository.GetAllUsers()
	if err != nil {
		t.Fatalf("expected all users to load: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected one user, got %d", len(users))
	}

	if nextID := repository.GetNextID(); nextID != 2 {
		t.Fatalf("expected next user ID 2, got %d", nextID)
	}
}

func useTempUserCSVConfig(t *testing.T) {
	t.Helper()

	previousPath := beego.AppConfig.DefaultString("csv_user_file", "data/users.csv")
	filePath := filepath.Join(t.TempDir(), "users.csv")

	if err := beego.AppConfig.Set("csv_user_file", filePath); err != nil {
		t.Fatalf("expected csv_user_file test config to be set: %v", err)
	}

	t.Cleanup(func() {
		_ = beego.AppConfig.Set("csv_user_file", previousPath)
	})
}
