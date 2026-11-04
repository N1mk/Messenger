package postgres

// Заглушка
func (r *FakePostgresRepo) ReplaceRefreshToken(newtokenHash []byte, userID int) error {
	return nil
}
