package repository

import "messenger/internal/domain"

// RefreshTokenExists checks if refresh token (domain.RefreshToken) with given tokenHash exists. Returns (true, nil) if refresh token exists, (false, nil) if not found and (false, err) if error occurred while check.
func (r *PostgresRepo) RefreshTokenExists(tokenHash []byte) (ok bool, err error) {
	var refreshToken domain.RefreshToken

	if err := r.db.Raw("SELECT * FROM refresh_tokens WHERE token_hash = ?", tokenHash).Scan(&refreshToken).Error; err != nil {
		return false, err
	}

	if refreshToken.TokenHash == nil {
		return false, nil
	}

	return true, nil
}
