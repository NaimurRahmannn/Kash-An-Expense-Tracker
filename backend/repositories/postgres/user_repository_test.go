package postgres

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"backend/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewUserRepository(t *testing.T) {
	repository, _, cleanup := newMockUserRepository(t)
	defer cleanup()

	if repository == nil {
		t.Fatal("expected repository to be non-nil")
	}
}

func TestUserRepositoryNilDBBehavior(t *testing.T) {
	repository := NewUserRepository(nil)

	tests := []struct {
		name string
		run  func() error
	}{
		{
			name: "GetAllUsers",
			run: func() error {
				_, err := repository.GetAllUsers()
				return err
			},
		},
		{
			name: "GetUserByEmail",
			run: func() error {
				_, err := repository.GetUserByEmail("john@example.com")
				return err
			},
		},
		{
			name: "GetUserByID",
			run: func() error {
				_, err := repository.GetUserByID(1)
				return err
			},
		},
		{
			name: "CreateUser",
			run: func() error {
				return repository.CreateUser(&models.User{Name: "John Doe"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if err == nil {
				t.Fatal("expected nil DB error")
			}
			if !strings.Contains(err.Error(), "database connection is nil") {
				t.Fatalf("expected nil DB error, got %q", err.Error())
			}
		})
	}

	if nextID := repository.GetNextID(); nextID != 1 {
		t.Fatalf("expected next ID 1 for nil DB, got %d", nextID)
	}
}

func TestUserRepositoryGetAllUsersReturnsUsers(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()
	createdAt := fixedPostgresTime()

	rows := sqlmock.NewRows(userColumns()).
		AddRow(1, "John Doe", "john@example.com", "secret123", createdAt).
		AddRow(2, "Jane Doe", "jane@example.com", "secret456", createdAt.Add(time.Hour))
	mock.ExpectQuery(regexp.QuoteMeta(selectAllUsersQuery)).WillReturnRows(rows)

	users, err := repository.GetAllUsers()
	if err != nil {
		t.Fatalf("expected users to load: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].ID != 1 || users[0].Email != "john@example.com" || users[0].CreatedAt != createdAt.UTC().Format(time.RFC3339) {
		t.Fatalf("unexpected first user: %+v", users[0])
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryGetAllUsersReturnsEmptySlice(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()

	rows := sqlmock.NewRows(userColumns())
	mock.ExpectQuery(regexp.QuoteMeta(selectAllUsersQuery)).WillReturnRows(rows)

	users, err := repository.GetAllUsers()
	if err != nil {
		t.Fatalf("expected users to load: %v", err)
	}
	if users == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(users) != 0 {
		t.Fatalf("expected no users, got %d", len(users))
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryGetAllUsersReturnsRowsError(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()

	rows := sqlmock.NewRows(userColumns()).
		AddRow(1, "John Doe", "john@example.com", "secret123", fixedPostgresTime()).
		RowError(0, errors.New("row scan failed"))
	mock.ExpectQuery(regexp.QuoteMeta(selectAllUsersQuery)).WillReturnRows(rows)

	if _, err := repository.GetAllUsers(); err == nil {
		t.Fatal("expected rows error")
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryGetUserByEmailReturnsUser(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()
	createdAt := fixedPostgresTime()

	rows := sqlmock.NewRows(userColumns()).
		AddRow(1, "John Doe", "john@example.com", "secret123", createdAt)
	mock.ExpectQuery(regexp.QuoteMeta(selectUserByEmailQuery)).
		WithArgs("john@example.com").
		WillReturnRows(rows)

	user, err := repository.GetUserByEmail(" john@example.com ")
	if err != nil {
		t.Fatalf("expected user lookup to succeed: %v", err)
	}
	if user == nil || user.ID != 1 || user.CreatedAt != createdAt.UTC().Format(time.RFC3339) {
		t.Fatalf("expected matching user, got %+v", user)
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryGetUserByEmailReturnsNilForMissingUser(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(selectUserByEmailQuery)).
		WithArgs("missing@example.com").
		WillReturnError(sql.ErrNoRows)

	user, err := repository.GetUserByEmail("missing@example.com")
	if err != nil {
		t.Fatalf("expected missing user to return nil error: %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryGetUserByIDReturnsUser(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()
	createdAt := fixedPostgresTime()

	rows := sqlmock.NewRows(userColumns()).
		AddRow(1, "John Doe", "john@example.com", "secret123", createdAt)
	mock.ExpectQuery(regexp.QuoteMeta(selectUserByIDQuery)).
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repository.GetUserByID(1)
	if err != nil {
		t.Fatalf("expected user lookup to succeed: %v", err)
	}
	if user == nil || user.ID != 1 || user.Email != "john@example.com" {
		t.Fatalf("expected matching user, got %+v", user)
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryGetUserByIDReturnsNilForMissingUser(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(selectUserByIDQuery)).
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	user, err := repository.GetUserByID(99)
	if err != nil {
		t.Fatalf("expected missing user to return nil error: %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryCreateUserAssignsReturnedID(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()
	user := &models.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	mock.ExpectQuery(regexp.QuoteMeta(insertUserQuery)).
		WithArgs(user.Name, user.Email, user.Password, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(7))

	if err := repository.CreateUser(user); err != nil {
		t.Fatalf("expected user to be created: %v", err)
	}
	if user.ID != 7 {
		t.Fatalf("expected returned ID 7, got %d", user.ID)
	}
	if user.CreatedAt == "" {
		t.Fatal("expected CreatedAt to be set")
	}
	if _, err := time.Parse(time.RFC3339, user.CreatedAt); err != nil {
		t.Fatalf("expected RFC3339 CreatedAt, got %q", user.CreatedAt)
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryCreateUserValidationErrors(t *testing.T) {
	repository, _, cleanup := newMockUserRepository(t)
	defer cleanup()

	if err := repository.CreateUser(nil); err == nil || !strings.Contains(err.Error(), "user is nil") {
		t.Fatalf("expected nil user error, got %v", err)
	}

	err := repository.CreateUser(&models.User{
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "secret123",
		CreatedAt: "invalid-date",
	})
	if err == nil {
		t.Fatal("expected invalid CreatedAt error")
	}
}

func TestUserRepositoryCreateUserReturnsDatabaseError(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()
	createdAt := fixedPostgresTime().UTC().Format(time.RFC3339)
	user := &models.User{
		Name:      "John Doe",
		Email:     "john@example.com",
		Password:  "secret123",
		CreatedAt: createdAt,
	}

	mock.ExpectQuery(regexp.QuoteMeta(insertUserQuery)).
		WithArgs(user.Name, user.Email, user.Password, sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	if err := repository.CreateUser(user); err == nil {
		t.Fatal("expected database error")
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryGetNextIDReturnsNextID(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(selectNextUserIDQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"next_id"}).AddRow(8))

	if nextID := repository.GetNextID(); nextID != 8 {
		t.Fatalf("expected next ID 8, got %d", nextID)
	}

	assertUserRepositoryExpectations(t, mock)
}

func TestUserRepositoryGetNextIDReturnsOneOnQueryError(t *testing.T) {
	repository, mock, cleanup := newMockUserRepository(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(selectNextUserIDQuery)).
		WillReturnError(errors.New("query failed"))

	if nextID := repository.GetNextID(); nextID != 1 {
		t.Fatalf("expected next ID 1, got %d", nextID)
	}

	assertUserRepositoryExpectations(t, mock)
}

func newMockUserRepository(t *testing.T) (*UserRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock database: %v", err)
	}

	return NewUserRepository(db), mock, func() {
		_ = db.Close()
	}
}

func assertUserRepositoryExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func fixedPostgresTime() time.Time {
	return time.Date(2025, 6, 1, 10, 30, 0, 0, time.UTC)
}

func userColumns() []string {
	return []string{"id", "name", "email", "password", "created_at"}
}
