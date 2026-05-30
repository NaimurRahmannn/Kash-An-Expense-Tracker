package postgres

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"backend/models"
)

const (
	selectAllUsersQuery    = "SELECT id, name, email, password, created_at FROM users ORDER BY id ASC"
	selectUserByEmailQuery = "SELECT id, name, email, password, created_at FROM users WHERE LOWER(email) = LOWER($1) LIMIT 1"
	selectUserByIDQuery    = "SELECT id, name, email, password, created_at FROM users WHERE id = $1 LIMIT 1"
	insertUserQuery        = "INSERT INTO users (name, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	selectNextUserIDQuery  = "SELECT COALESCE(MAX(id), 0) + 1 FROM users"
)

var (
	// ErrNilDB is returned when a repository method is called without a DB connection.
	ErrNilDB = errors.New("database connection is nil")
	// ErrNilUser is returned when CreateUser receives a nil user.
	ErrNilUser = errors.New("user is nil")
)

// UserRepository stores users in Postgres.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a Postgres user repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetAllUsers returns all users ordered by ID.
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	if r.db == nil {
		return nil, ErrNilDB
	}

	rows, err := r.db.Query(selectAllUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByEmail returns a user by case-insensitive email.
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	if r.db == nil {
		return nil, ErrNilDB
	}

	user, err := scanUser(r.db.QueryRow(selectUserByEmailQuery, strings.TrimSpace(email)))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return user, err
}

// GetUserByID returns a user by ID.
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	if r.db == nil {
		return nil, ErrNilDB
	}

	user, err := scanUser(r.db.QueryRow(selectUserByIDQuery, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return user, err
}

// CreateUser inserts a user and assigns the database-generated ID.
func (r *UserRepository) CreateUser(user *models.User) error {
	if r.db == nil {
		return ErrNilDB
	}
	if user == nil {
		return ErrNilUser
	}
	if user.CreatedAt == "" {
		user.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	createdAt, err := time.Parse(time.RFC3339, user.CreatedAt)
	if err != nil {
		return err
	}

	return r.db.QueryRow(
		insertUserQuery,
		user.Name,
		user.Email,
		user.Password,
		createdAt,
	).Scan(&user.ID)
}

// GetNextID returns the next available user ID for interface compatibility.
func (r *UserRepository) GetNextID() int {
	if r.db == nil {
		return 1
	}

	nextID := 1
	if err := r.db.QueryRow(selectNextUserIDQuery).Scan(&nextID); err != nil {
		return 1
	}

	return nextID
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (*models.User, error) {
	var user models.User
	var createdAt time.Time

	if err := scanner.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &createdAt); err != nil {
		return nil, err
	}

	user.CreatedAt = formatPostgresTime(createdAt)
	return &user, nil
}

func formatPostgresTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}
