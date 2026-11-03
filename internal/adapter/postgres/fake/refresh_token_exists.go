package postgres

// Заглушка
func (p *FakePostgres) RefreshTokenExists(tokenHash []byte) (ok bool, err error) {
	return true, nil
}
