package postgres

import (
	"database/sql"
	"strings"
	"testing"
)

func TestOpenDBRequiresDSN(t *testing.T) {
	_, err := OpenDB("  ")
	if err == nil {
		t.Fatal("expected empty DSN to return an error")
	}
	if !strings.Contains(err.Error(), "postgres dsn is required") {
		t.Fatalf("expected postgres DSN error, got %q", err.Error())
	}
}

func TestConfigurePoolDoesNotPanic(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://dummy")
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer db.Close()

	ConfigurePool(nil)
	ConfigurePool(db)
}

func TestRunMigrationsRequiresDB(t *testing.T) {
	err := RunMigrations(nil)
	if err == nil {
		t.Fatal("expected nil database to return an error")
	}
	if !strings.Contains(err.Error(), "database connection is nil") {
		t.Fatalf("expected nil database error, got %q", err.Error())
	}
}

func TestEmbeddedSchemaIsReadable(t *testing.T) {
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		t.Fatalf("expected embedded schema to be readable: %v", err)
	}

	content := string(schema)
	expectedSnippets := []string{
		"CREATE TABLE IF NOT EXISTS users",
		"CREATE TABLE IF NOT EXISTS expenses",
		"CREATE INDEX IF NOT EXISTS",
		"email TEXT NOT NULL UNIQUE",
		"user_id INTEGER NOT NULL",
		"amount NUMERIC(12,2) NOT NULL",
		"expense_date DATE NOT NULL",
	}

	for _, snippet := range expectedSnippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("expected schema to contain %q", snippet)
		}
	}
}
