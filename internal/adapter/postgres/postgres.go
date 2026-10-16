// Package postgres is an adapter for work with postgres databases.
package postgres

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/pkg/jwt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresInterface defines all methods for work with postgres databases.
type PostgresInterface interface {
	GetUserByID(id uint) (*domain.User, error)
	GetUserByName(username string) (*domain.User, error)
	CreateUser(user domain.User) (*domain.User, error)
	RemoveUser(id uint) error
	UpdateUser(user *domain.User) error

	GetChatByID(id uint) (*domain.Chat, error)
	GetChatsCount() int
	CreateChat(chat domain.Chat) (*domain.Chat, error)
	RemoveChat(id uint) error
	UpdateChat(chat *domain.Chat) error

	GetMessages(chatID uint) ([]domain.Message, error)

	RefreshTokenExists(token jwt.RefreshToken) (ok bool, err error)
	CreateRefreshToken(token jwt.RefreshToken) error
	RemoveRefreshToken(token jwt.RefreshToken) error
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

	db.AutoMigrate(&domain.User{}, &domain.Chat{}, &domain.Message{}, &jwt.RefreshToken{})

	return &Postgres{db: db}, nil
}
