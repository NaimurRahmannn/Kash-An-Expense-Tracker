package config

import (
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// StorageDriver identifies the configured storage backend.
type StorageDriver string

const (
	// StorageDriverCSV keeps CSV as the default assignment-compliant storage.
	StorageDriverCSV StorageDriver = "csv"
	// StorageDriverPostgres identifies the optional future production storage.
	StorageDriverPostgres StorageDriver = "postgres"
)

var (
	appConfigString = func(key string) (string, error) {
		return beego.AppConfig.String(key)
	}
	appConfigBool = func(key string) (bool, error) {
		return beego.AppConfig.Bool(key)
	}
)

// GetStorageDriver returns the configured storage driver, defaulting to CSV.
func GetStorageDriver() StorageDriver {
	value, err := appConfigString("storage_driver")
	if err != nil {
		return StorageDriverCSV
	}

	switch StorageDriver(strings.ToLower(strings.TrimSpace(value))) {
	case StorageDriverPostgres:
		return StorageDriverPostgres
	case StorageDriverCSV:
		return StorageDriverCSV
	default:
		return StorageDriverCSV
	}
}

// IsCSVStorage reports whether CSV storage is selected.
func IsCSVStorage() bool {
	return GetStorageDriver() == StorageDriverCSV
}

// IsPostgresStorage reports whether Postgres storage is selected.
func IsPostgresStorage() bool {
	return GetStorageDriver() == StorageDriverPostgres
}

// GetPostgresDSN returns the configured Postgres connection string.
func GetPostgresDSN() string {
	value, err := appConfigString("postgres_dsn")
	if err != nil {
		return ""
	}

	return strings.TrimSpace(value)
}

// IsPostgresAutoMigrateEnabled reports whether future Postgres auto-migration is enabled.
func IsPostgresAutoMigrateEnabled() bool {
	value, err := appConfigBool("postgres_auto_migrate")
	if err != nil {
		return false
	}

	return value
}
