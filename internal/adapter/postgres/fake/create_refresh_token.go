package postgres

// Заглушка
func (p *FakePostgres) CreateRefreshToken(token string, userID int) error {
	return nil
}
