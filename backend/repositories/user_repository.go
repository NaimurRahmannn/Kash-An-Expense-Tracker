package repositories

import "backend/models"

// UserRepository defines user persistence behavior across storage drivers.
type UserRepository interface {
	GetAllUsers() ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(user *models.User) error
	GetNextID() int
}
