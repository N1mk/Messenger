package postgres

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateChat updates chat with correct id.
func (p *Postgres) UpdateChat(chat *domain.Chat) error {
	if err := p.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
