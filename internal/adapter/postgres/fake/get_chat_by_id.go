package postgres

import "messenger/internal/domain"

func (p *FakePostgres) GetChatByID(id uint) (*domain.Chat, error) {
	return &p.Chats[id], nil
}
