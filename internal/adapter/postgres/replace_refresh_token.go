package postgres

import (
	"fmt"

	"gorm.io/gorm"
)

// ReplaceRefreshToken replaces the old user refresh token with the new one in transaction in the refresh_tokens table.
func (p *Postgres) ReplaceRefreshToken(newtoken string, userID int) error {
	if err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := p.RemoveRefreshToken(userID); err != nil {
			return fmt.Errorf("refresh token removing error: %w", err)
		}

		if err := p.CreateRefreshToken(newtoken, userID); err != nil {
			return fmt.Errorf("refresh token creating error: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
