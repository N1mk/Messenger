package postgres

import (
	"messenger/pkg/jwt"
)

// СreateRefreshToken creates new element in the refresh_tokens table.
func (p *Postgres) CreateRefreshToken(token jwt.RefreshToken) error {
	err := p.db.Create(&token).Error

	return err
}
