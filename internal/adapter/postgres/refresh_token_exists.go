package postgres

// RefreshTokenExists checks if refresh token exists in the refresh_tokens table.
func (p *Postgres) RefreshTokenExists(token string) (ok bool, err error) {
	var refreshToken RefreshToken

	if err := p.db.Raw("SELECT * FROM refresh_tokens WHERE token = ?", token).Scan(&refreshToken).Error; err != nil {
		return false, err
	}

	if refreshToken.Token == "" {
		return false, nil
	}

	return true, nil
}
