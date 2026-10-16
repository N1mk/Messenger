package postgres

import (
	"errors"
	"messenger/internal/domain"
	"messenger/pkg/jwt"

	"gorm.io/gorm"
)

// RemoveRefreshToken removes refresh token from refresh_tokens table.
func (p *Postgres) RemoveRefreshToken(token jwt.RefreshToken) error {
	if err := p.db.Exec("DELETE FROM refresh_tokens WHERE token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
