package postgres

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveChat removes chat from the chats table by id.
func (p *Postgres) RemoveChat(id int) error {
	if err := p.db.Exec("DELETE FROM user_chats WHERE chat_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	if err := p.db.Exec("DELETE FROM messages WHERE chat_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	if err := p.db.Exec("DELETE FROM chats WHERE id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
