package postgres

import "context"

// Заглушка
func (p *FakePostgres) RemoveRefreshToken(_ context.Context, userID int) error {
	return nil
}
