package postgres

import (
	"context"
	"messenger/internal/domain"
)

// GetChatsCount returns count of the chats table elements.
func (p *Postgres) GetChatsCount(ctx context.Context) int {
	var count int64

	p.db.WithContext(ctx).Model(&domain.Chat{}).Count(&count)

	return int(count)
}
