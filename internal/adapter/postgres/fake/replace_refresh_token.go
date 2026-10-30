package postgres

// Заглушка
func (p *FakePostgres) ReplaceRefreshToken(newtokenHash []byte, userID int) error {
	return nil
}
