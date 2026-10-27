// Package postgres is an adapter for work with postgres databases.
package postgres

import (
	"fmt"
	"messenger/internal/domain"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresInterface defines all methods for work with postgres databases.
type PostgresInterface interface {
	// GetUserByID selects user from the users table by id.
	GetUserByID(id uint) (*domain.User, error)

	// GetUserByName selects user from the users table by name.
	GetUserByName(username string) (*domain.User, error)

	// СreateUser creates new element in the users table.
	CreateUser(user domain.User) (*domain.User, error)

	// RemoveUser removes user from the users table by id.
	RemoveUser(id uint) error

	// UpdateUser updates user with correct id.
	UpdateUser(user *domain.User) error

	// GetChatByID selects chat from the chats table by id.
	GetChatByID(id uint) (*domain.Chat, error)

	// GetChatsCount returns count of the chats table elements.
	GetChatsCount() int

	// СreateChat creates new element in the chats table.
	CreateChat(chat domain.Chat) (*domain.Chat, error)

	// RemoveChat removes chat from the chats table by id.
	RemoveChat(id uint) error

	// UpdateChat updates chat with correct id.
	UpdateChat(chat *domain.Chat) error

	// GetMessages selects all messages from the messages table with correct chat id.
	GetMessages(chatID uint) ([]*domain.Message, error)

	// RefreshTokenExists checks if refresh token exists in the refresh_tokens table.
	RefreshTokenExists(token string) (ok bool, err error)

	// СreateRefreshToken creates new element in the refresh_tokens table.
	CreateRefreshToken(token string, userID uint) error

	// RemoveRefreshToken removes refresh token from refresh_tokens table.
	RemoveRefreshToken(userID uint) error
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
