package postgres

import (
	"messenger/pkg/jwt"
)

// Заглушка
func (p *FakePostgres) CreateRefreshToken(token jwt.RefreshToken) error {
	return nil
}
