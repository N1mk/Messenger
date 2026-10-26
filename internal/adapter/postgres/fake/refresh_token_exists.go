package postgres

// Заглушка
func (p *FakePostgres) RefreshTokenExists(token string) (ok bool, err error) {
	return true, nil
}
