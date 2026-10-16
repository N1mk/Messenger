package postgres

import (
	"messenger/pkg/jwt"
)

// Заглушка
func (p *FakePostgres) RemoveRefreshToken(token jwt.RefreshToken) error {
	return nil
}
