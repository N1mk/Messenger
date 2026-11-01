package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func RunMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose dialect setting error: %v", err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("migrations upping error: %v", err)
	}

	return nil
}
