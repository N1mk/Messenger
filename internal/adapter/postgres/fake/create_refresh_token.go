package postgres

// Заглушка
func (p *FakePostgres) CreateRefreshToken(token string, userID uint) error {
	return nil
}
