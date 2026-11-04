package fakerepository

// Заглушка
func (r *FakePostgresRepo) RefreshTokenExists(tokenHash []byte) (ok bool, err error) {
	return true, nil
}
