package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenDB opens and verifies a Postgres database connection.
func OpenDB(dsn string) (*sql.DB, error) {
	normalizedDSN := strings.TrimSpace(dsn)
	if normalizedDSN == "" {
		return nil, errors.New("postgres dsn is required")
	}

	db, err := sql.Open("pgx", normalizedDSN)
	if err != nil {
		return nil, err
	}

	ConfigurePool(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

// ConfigurePool applies default database connection pool settings.
func ConfigurePool(db *sql.DB) {
	if db == nil {
		return
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
}
