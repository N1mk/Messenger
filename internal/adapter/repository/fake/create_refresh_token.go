package postgres

// Заглушка
func (r *FakePostgresRepo) CreateRefreshToken(tokenHash []byte, userID int) error {
	return nil
}
