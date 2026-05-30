package storage

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	beego "github.com/beego/beego/v2/server/web"
)

func TestInitCSVModeDoesNotRequirePostgres(t *testing.T) {
	useStorageConfig(t, "csv", "", "false")
	if err := Close(); err != nil {
		t.Fatalf("expected close before test to succeed: %v", err)
	}

	if err := Init(); err != nil {
		t.Fatalf("expected CSV storage init to succeed: %v", err)
	}
	if PostgresDB() != nil {
		t.Fatal("expected CSV storage init to leave Postgres DB nil")
	}
	if err := Close(); err != nil {
		t.Fatalf("expected close after CSV init to succeed: %v", err)
	}
}

func TestInitPostgresModeRequiresDSN(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
	}{
		{
			name: "empty DSN",
		},
		{
			name: "whitespace DSN",
			dsn:  "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useStorageConfig(t, "postgres", tt.dsn, "false")
			if err := Close(); err != nil {
				t.Fatalf("expected close before test to succeed: %v", err)
			}

			err := Init()
			if err == nil {
				t.Fatal("expected missing Postgres DSN to return an error")
			}
			if !strings.Contains(err.Error(), "postgres dsn is required when storage_driver=postgres") {
				t.Fatalf("expected missing DSN error, got %q", err.Error())
			}
			if PostgresDB() != nil {
				t.Fatal("expected Postgres DB to remain nil when init fails")
			}
		})
	}
}

func TestInitPostgresModeWithInvalidDSNReturnsError(t *testing.T) {
	useStorageConfig(t, "postgres", "postgres://%", "false")
	if err := Close(); err != nil {
		t.Fatalf("expected close before test to succeed: %v", err)
	}

	if err := Init(); err == nil {
		t.Fatal("expected invalid Postgres DSN to return an error")
	}
	if PostgresDB() != nil {
		t.Fatal("expected Postgres DB to remain nil when init fails")
	}
}

func TestInitPostgresModeOpensDBWithoutMigration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock DB: %v", err)
	}
	mock.ExpectClose()

	useStorageConfig(t, "postgres", "postgres://example", "false")
	useStorageDependencies(
		t,
		func(dsn string) (*sql.DB, error) {
			if dsn != "postgres://example" {
				t.Fatalf("expected DSN %q, got %q", "postgres://example", dsn)
			}
			return db, nil
		},
		func(db *sql.DB) error {
			t.Fatal("expected migrations not to run when postgres_auto_migrate is false")
			return nil
		},
	)

	if err := Init(); err != nil {
		t.Fatalf("expected Postgres init to succeed: %v", err)
	}
	if PostgresDB() != db {
		t.Fatal("expected Postgres DB to be stored")
	}
	if err := Close(); err != nil {
		t.Fatalf("expected close to succeed: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected sqlmock expectations to be met: %v", err)
	}
}

func TestInitPostgresModeRunsMigrations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock DB: %v", err)
	}
	mock.ExpectClose()

	useStorageConfig(t, "postgres", "postgres://example", "true")
	migrationCalled := false
	useStorageDependencies(
		t,
		func(dsn string) (*sql.DB, error) {
			return db, nil
		},
		func(migrationDB *sql.DB) error {
			migrationCalled = true
			if migrationDB != db {
				t.Fatal("expected migration to receive opened DB")
			}
			return nil
		},
	)

	if err := Init(); err != nil {
		t.Fatalf("expected Postgres init with migration to succeed: %v", err)
	}
	if !migrationCalled {
		t.Fatal("expected migration to run")
	}
	if err := Close(); err != nil {
		t.Fatalf("expected close to succeed: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected sqlmock expectations to be met: %v", err)
	}
}

func TestInitPostgresModeClosesDBWhenMigrationFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock DB: %v", err)
	}
	mock.ExpectClose()

	useStorageConfig(t, "postgres", "postgres://example", "true")
	useStorageDependencies(
		t,
		func(dsn string) (*sql.DB, error) {
			return db, nil
		},
		func(db *sql.DB) error {
			return errors.New("migration failed")
		},
	)

	err = Init()
	if err == nil {
		t.Fatal("expected migration failure to be returned")
	}
	if !strings.Contains(err.Error(), "migration failed") {
		t.Fatalf("expected migration error, got %q", err.Error())
	}
	if PostgresDB() != nil {
		t.Fatal("expected Postgres DB to be reset after migration failure")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected sqlmock expectations to be met: %v", err)
	}
}

func TestInitPostgresModeReturnsOpenError(t *testing.T) {
	useStorageConfig(t, "postgres", "postgres://example", "false")
	useStorageDependencies(
		t,
		func(dsn string) (*sql.DB, error) {
			return nil, errors.New("open failed")
		},
		func(db *sql.DB) error {
			t.Fatal("expected migrations not to run when opening DB fails")
			return nil
		},
	)

	err := Init()
	if err == nil {
		t.Fatal("expected open error to be returned")
	}
	if !strings.Contains(err.Error(), "open failed") {
		t.Fatalf("expected open error, got %q", err.Error())
	}
	if PostgresDB() != nil {
		t.Fatal("expected Postgres DB to remain nil when open fails")
	}
}

func TestCloseCanBeCalledTwice(t *testing.T) {
	if err := Close(); err != nil {
		t.Fatalf("expected close before test to succeed: %v", err)
	}
	if err := Close(); err != nil {
		t.Fatalf("expected second close to succeed: %v", err)
	}
}

func TestCloseIsSafeWhenPostgresDBIsNil(t *testing.T) {
	if err := Close(); err != nil {
		t.Fatalf("expected close with nil Postgres DB to succeed: %v", err)
	}
}

func TestCloseResetsPostgresDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock DB: %v", err)
	}
	mock.ExpectClose()

	postgresDB = db
	t.Cleanup(func() {
		postgresDB = nil
	})

	if err := Close(); err != nil {
		t.Fatalf("expected close to succeed: %v", err)
	}
	if PostgresDB() != nil {
		t.Fatal("expected Postgres DB to be reset")
	}
	if err := Close(); err != nil {
		t.Fatalf("expected second close to succeed: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected sqlmock expectations to be met: %v", err)
	}
}

func useStorageConfig(t *testing.T, driver string, dsn string, autoMigrate string) {
	t.Helper()

	t.Setenv("STORAGE_DRIVER", "")
	t.Setenv("POSTGRES_DSN", "")
	t.Setenv("POSTGRES_AUTO_MIGRATE", "")

	previousDriver := beego.AppConfig.DefaultString("storage_driver", "csv")
	previousDSN := beego.AppConfig.DefaultString("postgres_dsn", "")
	previousAutoMigrate := beego.AppConfig.DefaultString("postgres_auto_migrate", "false")

	if err := beego.AppConfig.Set("storage_driver", driver); err != nil {
		t.Fatalf("expected storage_driver test config to be set: %v", err)
	}
	if err := beego.AppConfig.Set("postgres_dsn", dsn); err != nil {
		t.Fatalf("expected postgres_dsn test config to be set: %v", err)
	}
	if err := beego.AppConfig.Set("postgres_auto_migrate", autoMigrate); err != nil {
		t.Fatalf("expected postgres_auto_migrate test config to be set: %v", err)
	}

	t.Cleanup(func() {
		_ = Close()
		_ = beego.AppConfig.Set("storage_driver", previousDriver)
		_ = beego.AppConfig.Set("postgres_dsn", previousDSN)
		_ = beego.AppConfig.Set("postgres_auto_migrate", previousAutoMigrate)
	})
}

func useStorageDependencies(
	t *testing.T,
	openDB func(string) (*sql.DB, error),
	runMigrations func(*sql.DB) error,
) {
	t.Helper()

	previousOpenDB := openPostgresDB
	previousRunMigrations := runPostgresMigrations
	openPostgresDB = openDB
	runPostgresMigrations = runMigrations

	t.Cleanup(func() {
		openPostgresDB = previousOpenDB
		runPostgresMigrations = previousRunMigrations
	})
}
