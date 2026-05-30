package postgres

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRunMigrationsExecutesEmbeddedSchema(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock DB: %v", err)
	}
	defer db.Close()

	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		t.Fatalf("expected embedded schema to be readable: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(string(schema))).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := RunMigrations(db); err != nil {
		t.Fatalf("expected migrations to run: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected sqlmock expectations to be met: %v", err)
	}
}

func TestRunMigrationsReturnsExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("expected sqlmock DB: %v", err)
	}
	defer db.Close()

	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		t.Fatalf("expected embedded schema to be readable: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(string(schema))).
		WillReturnError(errors.New("migration failed"))

	err = RunMigrations(db)
	if err == nil {
		t.Fatal("expected migration error")
	}
	if err.Error() != "migration failed" {
		t.Fatalf("expected migration failure, got %q", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected sqlmock expectations to be met: %v", err)
	}
}
