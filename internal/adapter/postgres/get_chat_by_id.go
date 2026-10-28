package postgres

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetChatByID selects chat from the chats table by id.
func (p *Postgres) GetChatByID(id int) (*domain.Chat, error) {
	var chat domain.Chat

	err := p.db.Preload("Users").First(&chat, int(id)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &chat, nil
}
