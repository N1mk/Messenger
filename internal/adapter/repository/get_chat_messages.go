package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetChatMessages returns slice of pageSize messages (messages) starting starting with lastID (domain.Message) from chat with id = chatID and last message id (newLastID). If all message have been readen from chat returns -1 in newLastID. If chat with id = chatID not found GetChatMessages returns domain.ErrNotFound.
func (r *PostgresRepo) GetChatMessages(chatID int, lastID int, pageSize int) (messages []*domain.Message, newLastID int, err error) {
	err = r.db.Scopes(Paginate(lastID, pageSize)).Model(&domain.Message{}).Where("chat_id = ?", chatID).Select("*").Scan(&messages).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lastID, fmt.Errorf("chat not found: %w", domain.ErrNotFound)
		}
		return nil, lastID, err
	}

	return messages, lastID + len(messages), nil
}
