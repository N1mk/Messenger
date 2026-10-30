package postgres

// Заглушка
func (p *FakePostgres) CreateRefreshToken(tokenHash []byte, userID int) error {
	return nil
}
