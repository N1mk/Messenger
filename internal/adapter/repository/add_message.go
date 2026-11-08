package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// AddMessage adds new message (domain.Message) into the chat (domain.Chat) wiht id = message.ChatID. If chat with id = message.ChatID not found the AddMessage returns domain.ErrNotFound.
func (r *PostgresRepo) AddMessage(message *domain.Message) (*domain.Message, error) {
	if err := r.db.Model(&domain.Chat{}).First(message.ChatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("chat not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("chat exist check error: %w", err)
	}

	if err := r.db.Create(message).Error; err != nil {
		return nil, err
	}

	return message, nil
}
