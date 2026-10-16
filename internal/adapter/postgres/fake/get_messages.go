package postgres

import "messenger/internal/domain"

func (p *FakePostgres) GetMessages(chatID uint) ([]domain.Message, error) {
	return p.Chats[chatID].Messages, nil
}
