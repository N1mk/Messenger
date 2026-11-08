package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveRefreshToken removes refresh token (refresh tokens) (domain.RefreshToken) with given userID from the refresh_tokens table. If user with given id not foun RemoveRefreshToken returns domain.ErrNotFound.
func (r *PostgresRepo) RemoveRefreshToken(userID int) error {
	if err := r.db.Exec("DELETE FROM refresh_tokens WHERE user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return err
	}

	return nil
}
