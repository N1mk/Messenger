package postgres

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveUser removes user from the users table by id.
func (p *Postgres) RemoveUser(id int) error {
	if err := p.db.Exec("DELETE FROM user_chats WHERE user_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	if err := p.db.Delete(&domain.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
