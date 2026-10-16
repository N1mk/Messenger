package postgres

import (
	"messenger/pkg/jwt"
)

// Заглушка
func (p *FakePostgres) RefreshTokenExists(token jwt.RefreshToken) (ok bool, err error) {
	return true, nil
}
