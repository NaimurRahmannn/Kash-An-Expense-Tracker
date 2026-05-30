package repositories

import (
	"backend/config"
	csvrepo "backend/repositories/csv"
)

// NewUserRepository creates a user repository for the configured storage driver.
func NewUserRepository() UserRepository {
	switch config.GetStorageDriver() {
	case config.StorageDriverPostgres:
		// TODO: return Postgres repository when implemented.
		return csvrepo.NewUserRepository()
	case config.StorageDriverCSV:
		return csvrepo.NewUserRepository()
	default:
		return csvrepo.NewUserRepository()
	}
}

// NewExpenseRepository creates an expense repository for the configured storage driver.
func NewExpenseRepository() ExpenseRepository {
	switch config.GetStorageDriver() {
	case config.StorageDriverPostgres:
		// TODO: return Postgres repository when implemented.
		return csvrepo.NewExpenseRepository()
	case config.StorageDriverCSV:
		return csvrepo.NewExpenseRepository()
	default:
		return csvrepo.NewExpenseRepository()
	}
}
