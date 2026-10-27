package postgres

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetMessages selects all messages from the messages table with correct chat id.
func (p *Postgres) GetMessages(chatID uint) ([]*domain.Message, error) {
	var messages []domain.Message

	err := p.db.Raw("SELECT * FROM messages WHERE chat_id = ?", int(chatID)).Scan(&messages).Error
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
