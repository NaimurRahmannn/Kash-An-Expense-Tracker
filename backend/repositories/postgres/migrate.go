package postgres

import (
	"database/sql"
	"embed"
	"errors"
)

//go:embed schema.sql
var schemaFS embed.FS

// RunMigrations executes the embedded Postgres schema.
func RunMigrations(db *sql.DB) error {
	if db == nil {
		return errors.New("database connection is nil")
	}

	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schema))
	return err
}
