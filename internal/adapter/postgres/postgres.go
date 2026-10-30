// Package postgres is an adapter for work with postgres databases.
package postgres

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresInterface defines all methods for work with postgres databases.
type PostgresInterface interface {
	// GetUserByID selects user from the users table by id.
	GetUserByID(ctx context.Context, id int) (*domain.User, error)

	// GetUserByName selects user from the users table by name.
	GetUserByName(ctx context.Context, username string) (*domain.User, error)

	// СreateUser creates new element in the users table.
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)

	// RemoveUser removes user from the users table by id.
	RemoveUser(ctx context.Context, id int) error

	// UpdateUser updates user with correct id.
	UpdateUser(ctx context.Context, user *domain.User) error

	// GetChatByID selects chat from the chats table by id.
	GetChatByID(ctx context.Context, id int) (*domain.Chat, error)

	// GetChatsCount returns count of the chats table elements.
	GetChatsCount(ctx context.Context) int

	// СreateChat creates new element in the chats table.
	CreateChat(ctx context.Context, chat domain.Chat) (*domain.Chat, error)

	// RemoveChat removes chat from the chats table by id.
	RemoveChat(ctx context.Context, id int) error

	// UpdateChat updates chat with correct id.
	UpdateChat(ctx context.Context, chat *domain.Chat) error

	// GetMessages selects all messages from the messages table with correct chat id.
	GetMessages(ctx context.Context, chatID int) ([]*domain.Message, error)

	// RefreshTokenExists checks if refresh token exists in the refresh_tokens table.
	RefreshTokenExists(ctx context.Context, tokenHash []byte) (ok bool, err error)

	// СreateRefreshToken creates new element in the refresh_tokens table.
	CreateRefreshToken(ctx context.Context, tokenHash []byte, userID int) error

	// RemoveRefreshToken removes refresh token from refresh_tokens table.
	RemoveRefreshToken(ctx context.Context, userID int) error

	// ReplaceRefreshToken replaces the old user refresh token with the new one in transaction in the refresh_tokens table.
	ReplaceRefreshToken(ctx context.Context, newtokenHash []byte, userID int) error
}

// Postgres implements PostgresInterface and containts *gorm.DB for work with postgres databases.
type Postgres struct {
	db *gorm.DB
}

// NewPostgres initializes Postgres.
// The function returns pointer on Postgres struct and error.
func NewPostgres(dsn string) (*Postgres, error) {
	if dsn == "" {
		dsn = os.Getenv("DEFAULT_DSN")
		if dsn == "" {
			return nil, fmt.Errorf("no default dsn in env variables")
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&domain.User{}, &domain.Chat{}, &domain.Message{}, &RefreshToken{})

	return &Postgres{db: db}, nil
}
