package postgres

import "context"

// Заглушка
func (p *FakePostgres) CreateRefreshToken(_ context.Context, tokenHash []byte, userID int) error {
	return nil
}
