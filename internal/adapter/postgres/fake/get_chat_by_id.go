package postgres

import "messenger/internal/domain"

func (p *FakePostgres) GetChatByID(id int) (*domain.Chat, error) {
	return &p.Chats[id], nil
}
