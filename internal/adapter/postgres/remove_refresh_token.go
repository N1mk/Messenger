package postgres

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveRefreshToken removes refresh token from refresh_tokens table.
func (p *Postgres) RemoveRefreshToken(userID int) error {
	if err := p.db.Exec("DELETE FROM refresh_tokens WHERE user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
