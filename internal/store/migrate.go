package store

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB, migrationsDir string) error {
	absDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("resolve migrations dir: %w", err)
	}
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.Up(db, absDir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}
