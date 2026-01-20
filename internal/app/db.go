package app

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"log"

	_ "github.com/mattn/go-sqlite3"
	goose "github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func OpenDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	goose.SetLogger(log.New(io.Discard, "", 0))

	goose.SetBaseFS(migrationsFS)

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
