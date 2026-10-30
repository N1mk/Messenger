package postgres

import "context"

// RefreshTokenExists checks if refresh token exists in the refresh_tokens table.
func (p *Postgres) RefreshTokenExists(ctx context.Context, tokenHash []byte) (ok bool, err error) {
	var refreshToken RefreshToken

	if err := p.db.WithContext(ctx).Raw("SELECT * FROM refresh_tokens WHERE token_hash = ?", tokenHash).Scan(&refreshToken).Error; err != nil {
		return false, err
	}

	if refreshToken.TokenHash == nil {
		return false, nil
	}

	return true, nil
}
