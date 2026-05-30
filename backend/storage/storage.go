package storage

import (
	"database/sql"
	"errors"

	"backend/config"
	"backend/repositories/postgres"
)

var postgresDB *sql.DB

// Init initializes the configured storage backend.
func Init() error {
	if config.GetStorageDriver() == config.StorageDriverCSV {
		return nil
	}

	dsn := config.GetPostgresDSN()
	if dsn == "" {
		return errors.New("postgres dsn is required when storage_driver=postgres")
	}

	db, err := postgres.OpenDB(dsn)
	if err != nil {
		return err
	}
	postgresDB = db

	if config.IsPostgresAutoMigrateEnabled() {
		if err := postgres.RunMigrations(postgresDB); err != nil {
			_ = Close()
			return err
		}
	}

	return nil
}

// Close closes the shared Postgres connection if it exists.
func Close() error {
	if postgresDB == nil {
		return nil
	}

	err := postgresDB.Close()
	postgresDB = nil
	return err
}

// PostgresDB returns the shared Postgres connection.
func PostgresDB() *sql.DB {
	return postgresDB
}
