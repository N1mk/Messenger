package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetChatByID recives user's id and returns *domain.Chat with same id or domain.ErrNotFound if user is not found.
func (r *PostgresRepo) GetChatByID(id int) (*domain.Chat, error) {
	var chat domain.Chat

	err := r.db.Preload("Users").First(&chat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("chat not found: %w", domain.ErrNotFound)
		}
		return nil, err
	}

	return &chat, nil
}
