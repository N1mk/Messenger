package repository

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetMessages selects all messages from the messages table with correct chat id.
func (r *PostgresRepo) GetMessages(chatID int, lastID int, pageSize int) ([]*domain.Message, error) {
	var messages []domain.Message

	err := r.db.Scopes(Paginate(lastID, pageSize)).Model(&domain.Message{}).Where("chat_id = ?", chatID).Select("*").Scan(&messages).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	messagePointers := make([]*domain.Message, len(messages))
	for i, message := range messages {
		messagePointers[i] = &message
	}

	return messagePointers, nil
}
