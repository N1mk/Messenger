package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetChatMessages returns slice of pageSize messages starting starting with lastID (domain.Message) from chat with id = chatID. If chat with id = chatID not found GetChatMessages returns domain.ErrNotFound.
func (r *PostgresRepo) GetChatMessages(chatID int, lastID int, pageSize int) ([]*domain.Message, error) {
	var messages []*domain.Message

	err := r.db.Scopes(Paginate(lastID, pageSize)).Model(&domain.Message{}).Where("chat_id = ?", chatID).Select("*").Scan(&messages).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("chat not found: %w", domain.ErrNotFound)
		}
		return nil, err
	}

	return messages, nil
}
