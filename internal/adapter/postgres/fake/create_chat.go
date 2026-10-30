package postgres

import (
	"context"
	"messenger/internal/domain"
)

func (p *FakePostgres) CreateChat(_ context.Context, chat domain.Chat) (*domain.Chat, error) {
	chat.ID = len(p.Chats)

	p.Chats = append(p.Chats, chat)

	return &chat, nil
}
