package postgres

import (
	"context"
	"messenger/internal/domain"
)

func (p *FakePostgres) GetChatByID(_ context.Context, id int) (*domain.Chat, error) {
	return &p.Chats[id], nil
}
