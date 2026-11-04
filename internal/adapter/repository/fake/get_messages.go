package postgres

import (
	"messenger/internal/domain"
)

func (r *FakePostgresRepo) GetMessages(chatID int) ([]*domain.Message, error) {
	messages := r.Chats[chatID].Messages

	messagePointers := make([]*domain.Message, len(messages))
	for i, message := range messages {
		messagePointers[i] = &message
	}

	return messagePointers, nil
}
