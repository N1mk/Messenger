package postgres

import "messenger/internal/domain"

// СreateChat creates new element in the chats table.
func (p *Postgres) CreateChat(chat domain.Chat) (*domain.Chat, error) {
	err := p.db.Create(&chat).Error

	if err != nil {
		return nil, err
	}

	return &chat, nil
}
