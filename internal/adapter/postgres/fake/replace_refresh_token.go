package postgres

// Заглушка
func (p *FakePostgres) ReplaceRefreshToken(newtoken string, userID int) error {
	return nil
}
