package repositories

import (
	"backend/config"
	csvrepo "backend/repositories/csv"
	postgresrepo "backend/repositories/postgres"
	"backend/storage"
)

var getStorageDriver = config.GetStorageDriver

// NewUserRepository creates a user repository for the configured storage driver.
func NewUserRepository() UserRepository {
	switch getStorageDriver() {
	case config.StorageDriverPostgres:
		return postgresrepo.NewUserRepository(storage.PostgresDB())
	case config.StorageDriverCSV:
		return csvrepo.NewUserRepository()
	default:
		return csvrepo.NewUserRepository()
	}
}

// NewExpenseRepository creates an expense repository for the configured storage driver.
func NewExpenseRepository() ExpenseRepository {
	switch getStorageDriver() {
	case config.StorageDriverPostgres:
		return postgresrepo.NewExpenseRepository(storage.PostgresDB())
	case config.StorageDriverCSV:
		return csvrepo.NewExpenseRepository()
	default:
		return csvrepo.NewExpenseRepository()
	}
}
