package postgres

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateUser updates user with correct id.
func (p *Postgres) UpdateUser(user *domain.User) error {
	if err := p.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
