package config

import (
	"errors"
	"testing"
)

func TestGetStorageDriver(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		readError bool
		want      StorageDriver
	}{
		{
			name:      "missing storage_driver defaults to csv",
			readError: true,
			want:      StorageDriverCSV,
		},
		{
			name:  "empty storage_driver defaults to csv",
			value: "   ",
			want:  StorageDriverCSV,
		},
		{
			name:  "csv returns csv",
			value: "csv",
			want:  StorageDriverCSV,
		},
		{
			name:  "postgres returns postgres",
			value: "postgres",
			want:  StorageDriverPostgres,
		},
		{
			name:  "postgres trims spaces and ignores case",
			value: " Postgres ",
			want:  StorageDriverPostgres,
		},
		{
			name:  "unknown storage_driver defaults to csv",
			value: "sqlite",
			want:  StorageDriverCSV,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearStorageEnvironment(t)
			useTestStringConfig(t, map[string]string{"storage_driver": tt.value}, tt.readError)

			if got := GetStorageDriver(); got != tt.want {
				t.Fatalf("expected storage driver %q, got %q", tt.want, got)
			}
		})
	}
}

func TestStorageDriverChecks(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		wantCSV      bool
		wantPostgres bool
	}{
		{
			name:    "csv storage",
			value:   "csv",
			wantCSV: true,
		},
		{
			name:         "postgres storage",
			value:        "postgres",
			wantPostgres: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearStorageEnvironment(t)
			useTestStringConfig(t, map[string]string{"storage_driver": tt.value}, false)

			if got := IsCSVStorage(); got != tt.wantCSV {
				t.Fatalf("expected IsCSVStorage %v, got %v", tt.wantCSV, got)
			}
			if got := IsPostgresStorage(); got != tt.wantPostgres {
				t.Fatalf("expected IsPostgresStorage %v, got %v", tt.wantPostgres, got)
			}
		})
	}
}

func TestGetPostgresDSN(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		readError bool
		want      string
	}{
		{
			name: "empty postgres_dsn returns empty string",
			want: "",
		},
		{
			name:      "missing postgres_dsn returns empty string",
			readError: true,
			want:      "",
		},
		{
			name:  "postgres_dsn trims spaces",
			value: " postgres://user:pass@localhost:5432/expenses ",
			want:  "postgres://user:pass@localhost:5432/expenses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearStorageEnvironment(t)
			useTestStringConfig(t, map[string]string{"postgres_dsn": tt.value}, tt.readError)

			if got := GetPostgresDSN(); got != tt.want {
				t.Fatalf("expected postgres DSN %q, got %q", tt.want, got)
			}
		})
	}
}

func TestIsPostgresAutoMigrateEnabled(t *testing.T) {
	tests := []struct {
		name      string
		value     bool
		readError bool
		want      bool
	}{
		{
			name:      "postgres_auto_migrate false by default",
			readError: true,
			want:      false,
		},
		{
			name:  "postgres_auto_migrate false returns false",
			value: false,
			want:  false,
		},
		{
			name:  "postgres_auto_migrate true returns true",
			value: true,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearStorageEnvironment(t)
			useTestBoolConfig(t, tt.value, tt.readError)

			if got := IsPostgresAutoMigrateEnabled(); got != tt.want {
				t.Fatalf("expected auto migrate %v, got %v", tt.want, got)
			}
		})
	}
}

func TestStorageEnvironmentOverrides(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "storage driver environment overrides app config",
			run: func(t *testing.T) {
				useTestStringConfig(t, map[string]string{"storage_driver": "csv"}, false)
				t.Setenv("STORAGE_DRIVER", "postgres")

				if got := GetStorageDriver(); got != StorageDriverPostgres {
					t.Fatalf("expected storage driver %q, got %q", StorageDriverPostgres, got)
				}
			},
		},
		{
			name: "postgres dsn environment overrides app config",
			run: func(t *testing.T) {
				useTestStringConfig(t, map[string]string{"postgres_dsn": ""}, false)
				t.Setenv("POSTGRES_DSN", " postgres://user:pass@host:5432/db ")

				want := "postgres://user:pass@host:5432/db"
				if got := GetPostgresDSN(); got != want {
					t.Fatalf("expected postgres DSN %q, got %q", want, got)
				}
			},
		},
		{
			name: "postgres auto migrate environment overrides app config",
			run: func(t *testing.T) {
				useTestBoolConfig(t, false, false)
				t.Setenv("POSTGRES_AUTO_MIGRATE", "true")

				if !IsPostgresAutoMigrateEnabled() {
					t.Fatal("expected postgres auto migrate to be enabled")
				}
			},
		},
		{
			name: "invalid postgres auto migrate environment is false",
			run: func(t *testing.T) {
				useTestBoolConfig(t, true, false)
				t.Setenv("POSTGRES_AUTO_MIGRATE", "not-bool")

				if IsPostgresAutoMigrateEnabled() {
					t.Fatal("expected invalid postgres auto migrate environment to be false")
				}
			},
		},
		{
			name: "csv user file environment overrides app config",
			run: func(t *testing.T) {
				useTestStringConfig(t, map[string]string{"csv_user_file": "data/users.csv"}, false)
				t.Setenv("CSV_USER_FILE", " /app/data/users.csv ")

				if got := GetCSVUserFile(); got != "/app/data/users.csv" {
					t.Fatalf("expected CSV user file %q, got %q", "/app/data/users.csv", got)
				}
			},
		},
		{
			name: "csv expense file environment overrides app config",
			run: func(t *testing.T) {
				useTestStringConfig(t, map[string]string{"csv_expense_file": "data/expenses.csv"}, false)
				t.Setenv("CSV_EXPENSE_FILE", " /app/data/expenses.csv ")

				if got := GetCSVExpenseFile(); got != "/app/data/expenses.csv" {
					t.Fatalf("expected CSV expense file %q, got %q", "/app/data/expenses.csv", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearStorageEnvironment(t)
			tt.run(t)
		})
	}
}

func clearStorageEnvironment(t *testing.T) {
	t.Helper()

	t.Setenv("STORAGE_DRIVER", "")
	t.Setenv("POSTGRES_DSN", "")
	t.Setenv("POSTGRES_AUTO_MIGRATE", "")
	t.Setenv("CSV_USER_FILE", "")
	t.Setenv("CSV_EXPENSE_FILE", "")
}

func useTestStringConfig(t *testing.T, values map[string]string, readError bool) {
	t.Helper()

	previous := appConfigString
	appConfigString = func(key string) (string, error) {
		if readError {
			return "", errors.New("missing config")
		}

		return values[key], nil
	}

	t.Cleanup(func() {
		appConfigString = previous
	})
}

func useTestBoolConfig(t *testing.T, value bool, readError bool) {
	t.Helper()

	previous := appConfigBool
	appConfigBool = func(key string) (bool, error) {
		if readError {
			return false, errors.New("missing config")
		}

		return value, nil
	}

	t.Cleanup(func() {
		appConfigBool = previous
	})
}
