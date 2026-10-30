package postgres

import "context"

// Заглушка
func (p *FakePostgres) RefreshTokenExists(_ context.Context, tokenHash []byte) (ok bool, err error) {
	return true, nil
}
