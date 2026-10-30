package postgres

import (
	"context"
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetChatByID selects chat from the chats table by id.
func (p *Postgres) GetChatByID(ctx context.Context, id int) (*domain.Chat, error) {
	var chat domain.Chat

	err := p.db.WithContext(ctx).Preload("Users").First(&chat, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &chat, nil
}
