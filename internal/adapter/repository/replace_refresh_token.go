package postgres

import (
	"fmt"
)

// ReplaceRefreshToken replaces the old user refresh token with the new one in transaction in the refresh_tokens table.
func (r *PostgresRepo) ReplaceRefreshToken(newtokenHash []byte, userID int) error {
	if err := r.RemoveRefreshToken(userID); err != nil {
		return fmt.Errorf("refresh token removing error: %w", err)
	}

	if err := r.CreateRefreshToken(newtokenHash, userID); err != nil {
		return fmt.Errorf("refresh token creating error: %w", err)
	}

	return nil
}
