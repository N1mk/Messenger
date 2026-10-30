package postgres

import (
	"context"
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveUser removes user from the users table by id.
func (p *Postgres) RemoveUser(ctx context.Context, id int) error {
	if err := p.db.WithContext(ctx).Exec("DELETE FROM user_chats WHERE user_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	if err := p.db.WithContext(ctx).Delete(&domain.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
