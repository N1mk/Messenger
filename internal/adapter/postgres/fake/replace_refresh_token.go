package postgres

import "context"

// Заглушка
func (p *FakePostgres) ReplaceRefreshToken(_ context.Context, newtokenHash []byte, userID int) error {
	return nil
}
