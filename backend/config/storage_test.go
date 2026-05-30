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
			useTestBoolConfig(t, tt.value, tt.readError)

			if got := IsPostgresAutoMigrateEnabled(); got != tt.want {
				t.Fatalf("expected auto migrate %v, got %v", tt.want, got)
			}
		})
	}
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
