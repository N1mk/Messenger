package postgres

import (
	"context"
	"messenger/internal/domain"
)

func (p *FakePostgres) UpdateChat(_ context.Context, chat *domain.Chat) error {
	p.Chats[chat.ID] = *chat

	return nil
}
