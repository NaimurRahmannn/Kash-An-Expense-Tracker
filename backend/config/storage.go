package config

import (
	"os"
	"strconv"
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
	value := getStringConfig("storage_driver", "STORAGE_DRIVER")

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
	return getStringConfig("postgres_dsn", "POSTGRES_DSN")
}

// GetCSVUserFile returns the configured CSV user file path.
func GetCSVUserFile() string {
	return getStringConfig("csv_user_file", "CSV_USER_FILE")
}

// GetCSVExpenseFile returns the configured CSV expense file path.
func GetCSVExpenseFile() string {
	return getStringConfig("csv_expense_file", "CSV_EXPENSE_FILE")
}

// IsPostgresAutoMigrateEnabled reports whether future Postgres auto-migration is enabled.
func IsPostgresAutoMigrateEnabled() bool {
	envValue := strings.TrimSpace(os.Getenv("POSTGRES_AUTO_MIGRATE"))
	if envValue != "" {
		enabled, err := strconv.ParseBool(envValue)
		return err == nil && enabled
	}

	value, err := appConfigBool("postgres_auto_migrate")
	if err != nil {
		return false
	}

	return value
}

func getStringConfig(configKey string, envKey string) string {
	if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
		return value
	}

	value, err := appConfigString(configKey)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(value)
}
