// Package postgres is an adapter for work with postgres databases.
package postgres

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/migrations"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresInterface defines all methods for work with postgres databases.
type PostgresInterface interface {
	// Begins transaction
	Begin(ctx context.Context)
	// Rollbacks transaction
	Rollback()
	// Commits transaction
	Commit()
	// GetUserByID selects user from the users table by id.
	GetUserByID(id int) (*domain.User, error)

	// GetUserByName selects user from the users table by name.
	GetUserByName(username string) (*domain.User, error)

	// СreateUser creates new element in the users table.
	CreateUser(user domain.User) (*domain.User, error)

	// RemoveUser removes user from the users table by id.
	RemoveUser(id int) error

	// UpdateUser updates user with correct id.
	UpdateUser(user *domain.User) error

	// GetChatByID selects chat from the chats table by id.
	GetChatByID(id int) (*domain.Chat, error)

	// GetChatsCount returns count of the chats table elements.
	GetChatsCount() int

	// СreateChat creates new element in the chats table.
	CreateChat(chat domain.Chat) (*domain.Chat, error)

	// RemoveChat removes chat from the chats table by id.
	RemoveChat(id int) error

	// UpdateChat updates chat with correct id.
	UpdateChat(chat *domain.Chat) error

	// GetMessages selects all messages from the messages table with correct chat id.
	GetMessages(chatID int) ([]*domain.Message, error)

	// RefreshTokenExists checks if refresh token exists in the refresh_tokens table.
	RefreshTokenExists(tokenHash []byte) (ok bool, err error)

	// СreateRefreshToken creates new element in the refresh_tokens table.
	CreateRefreshToken(tokenHash []byte, userID int) error

	// RemoveRefreshToken removes refresh token from refresh_tokens table.
	RemoveRefreshToken(userID int) error

	// ReplaceRefreshToken replaces the old user refresh token with the new one in transaction in the refresh_tokens table.
	ReplaceRefreshToken(newtokenHash []byte, userID int) error
}

// Postgres implements PostgresInterface and containts *gorm.DB for work with postgres databases.
type Postgres struct {
	db    *gorm.DB
	oldDB *gorm.DB
}

// NewPostgres initializes Postgres.
// The function returns pointer on Postgres struct and error.
func NewPostgres(dsn string) (*Postgres, error) {
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

	return &Postgres{db: db}, nil
}

// Saves db into oldDB and begins transaction in db
func (p *Postgres) Begin(ctx context.Context) {
	currentDB := *p.db
	p.oldDB = &currentDB

	p.db = p.db.WithContext(ctx).Begin()
}

// Rollbacks transaction and changes db to oldDB
func (p *Postgres) Rollback() {
	p.db.Rollback()

	oldDB := *p.oldDB
	p.db = &oldDB
}

// Rollbacks transaction and changes db to oldDB
func (p *Postgres) Commit() {
	p.db.Commit()

	oldDB := *p.oldDB
	p.db = &oldDB
}
