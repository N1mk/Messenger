package repository

import "messenger/internal/domain"

// AddRefreshToken adds new refresh token (domain.RefreshToken) into the refresh_tokens table.
func (r *PostgresRepo) AddRefreshToken(token *domain.RefreshToken) (*domain.RefreshToken, error) {
	if err := r.db.Model(&domain.RefreshToken{}).Create(token).Error; err != nil {
		return nil, err
	}

	return token, nil
}
