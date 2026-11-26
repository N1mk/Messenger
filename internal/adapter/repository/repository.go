// Package postgres is an adapter for work with postgres databases.
package repository

import (
	"context"
	"fmt"
	"messenger/internal/usecase"
	"messenger/migrations"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresRepo implements the Repo interface and wroks with postgres databases.
type PostgresRepo struct {
	db    *gorm.DB
	oldDB *gorm.DB
	mb    usecase.MessageBroker
}

// NewPostgres initializes Postgres.
// The function returns pointer on Postgres struct and error.
func NewPostgres(dsn string, mb usecase.MessageBroker) (*PostgresRepo, error) {
	if dsn == "" {
		dsn = os.Getenv("DEFAULT_DSN")
		if dsn == "" {
			return nil, fmt.Errorf("default dsn not found in env variables")
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("db opening error: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("sql.DB getting error: %w", err)
	}

	if err := migrations.RunMigrations(sqlDB); err != nil {
		return nil, fmt.Errorf("migrations running error: %w", err)
	}

	return &PostgresRepo{db: db, mb: mb}, nil
}

// Saves db into oldDB and begins transaction in db
func (pr *PostgresRepo) Begin(ctx context.Context) {
	currentDB := *pr.db
	pr.oldDB = &currentDB

	pr.db = pr.db.WithContext(ctx).Begin()
}

// Rollbacks transaction and changes db to oldDB
func (pr *PostgresRepo) Rollback() {
	pr.db.Rollback()

	oldDB := *pr.oldDB
	pr.db = &oldDB
}

// Rollbacks transaction and changes db to oldDB
func (pr *PostgresRepo) Commit() {
	pr.db.Commit()

	oldDB := *pr.oldDB
	pr.db = &oldDB
}
