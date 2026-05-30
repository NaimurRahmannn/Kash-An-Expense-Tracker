package storage

import (
	"strings"
	"testing"

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
	useStorageConfig(t, "postgres", "   ", "false")
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
}

func TestCloseIsSafeWhenPostgresDBIsNil(t *testing.T) {
	if err := Close(); err != nil {
		t.Fatalf("expected close with nil Postgres DB to succeed: %v", err)
	}
}

func useStorageConfig(t *testing.T, driver string, dsn string, autoMigrate string) {
	t.Helper()

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
