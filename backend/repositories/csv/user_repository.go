package csv

import "backend/models"

// UserRepository adapts the existing CSV-backed user model functions.
type UserRepository struct{}

// NewUserRepository creates a CSV user repository adapter.
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// GetAllUsers returns all users from CSV storage.
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	return models.GetAllUsers()
}

// GetUserByEmail returns a user by email from CSV storage.
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	return models.GetUserByEmail(email)
}

// GetUserByID returns a user by ID from CSV storage.
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	return models.GetUserByID(id)
}

// CreateUser stores a user in CSV storage.
func (r *UserRepository) CreateUser(user *models.User) error {
	return models.CreateUser(user)
}

// GetNextID returns the next available CSV user ID.
func (r *UserRepository) GetNextID() int {
	return models.GetNextID()
}
