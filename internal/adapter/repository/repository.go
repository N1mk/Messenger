// Package postgres is an adapter for work with postgres databases.
package repository

import (
	"context"
	"fmt"
	"messenger/migrations"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresRepo implements PostgresInterface and containts *gorm.DB for work with postgres databases.
type PostgresRepo struct {
	db    *gorm.DB
	oldDB *gorm.DB
}

// NewPostgres initializes Postgres.
// The function returns pointer on Postgres struct and error.
func NewPostgres(dsn string) (*PostgresRepo, error) {
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

	return &PostgresRepo{db: db}, nil
}

// Saves db into oldDB and begins transaction in db
func (r *PostgresRepo) Begin(ctx context.Context) {
	currentDB := *r.db
	r.oldDB = &currentDB

	r.db = r.db.WithContext(ctx).Begin()
}

// Rollbacks transaction and changes db to oldDB
func (r *PostgresRepo) Rollback() {
	r.db.Rollback()

	oldDB := *r.oldDB
	r.db = &oldDB
}

// Rollbacks transaction and changes db to oldDB
func (r *PostgresRepo) Commit() {
	r.db.Commit()

	oldDB := *r.oldDB
	r.db = &oldDB
}
