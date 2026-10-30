package postgres

import (
	"context"
	"messenger/internal/domain"
)

// СreateChat creates new element in the chats table.
func (p *Postgres) CreateChat(ctx context.Context, chat domain.Chat) (*domain.Chat, error) {
	err := p.db.WithContext(ctx).Create(&chat).Error

	if err != nil {
		return nil, err
	}

	return &chat, nil
}
