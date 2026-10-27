package postgres

import "messenger/internal/domain"

func (p *FakePostgres) GetMessages(chatID uint) ([]*domain.Message, error) {
	messages := p.Chats[chatID].Messages

	messagePointers := make([]*domain.Message, len(messages))
	for i, message := range messages {
		messagePointers[i] = &message
	}

	return messagePointers, nil
}
