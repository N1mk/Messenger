package postgres

import "messenger/internal/domain"

// GetChatsCount returns count of the chats table elements.
func (p *Postgres) GetChatsCount() int {
	var count int64

	p.db.Model(&domain.Chat{}).Count(&count)

	return int(count)
}
