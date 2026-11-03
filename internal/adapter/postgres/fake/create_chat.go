package postgres

import (
	"messenger/internal/domain"
)

func (p *FakePostgres) CreateChat(chat domain.Chat) (*domain.Chat, error) {
	chat.ID = len(p.Chats)

	p.Chats = append(p.Chats, chat)

	return &chat, nil
}
