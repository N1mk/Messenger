package repository

import (
	"fmt"
	"messenger/internal/domain"
)

// ReplaceRefreshToken replaces refresh token of user with id = userID. If user with id = userID not found ReplaceRefreshToken returns domain.NotFound.
func (r *PostgresRepo) ReplaceRefreshToken(newToken *domain.RefreshToken) error {
	if err := r.RemoveRefreshToken(newToken.UserID); err != nil {
		return fmt.Errorf("refresh token removing error: %w", err)
	}

	if _, err := r.AddRefreshToken(newToken); err != nil {
		return fmt.Errorf("refresh token creating error: %w", err)
	}

	return nil
}
