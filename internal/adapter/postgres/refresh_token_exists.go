package postgres

// RefreshTokenExists checks if refresh token exists in the refresh_tokens table.
func (p *Postgres) RefreshTokenExists(tokenHash []byte) (ok bool, err error) {
	var refreshToken RefreshToken

	if err := p.db.Raw("SELECT * FROM refresh_tokens WHERE token_hash = ?", tokenHash).Scan(&refreshToken).Error; err != nil {
		return false, err
	}

	return true, nil
}
