package postgres

import (
	"errors"
	"messenger/pkg/jwt"

	"gorm.io/gorm"
)

// RefreshTokenExists checks if refresh token exists in the refresh_tokens table.
func (p *Postgres) RefreshTokenExists(token jwt.RefreshToken) (ok bool, err error) {
	if err := p.db.Raw("SELECT * FROM refresh_tokens WHERE token = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
