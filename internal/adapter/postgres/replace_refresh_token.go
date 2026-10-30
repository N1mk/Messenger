package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// ReplaceRefreshToken replaces the old user refresh token with the new one in transaction in the refresh_tokens table.
func (p *Postgres) ReplaceRefreshToken(ctx context.Context, newtokenHash []byte, userID int) error {
	if err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := p.RemoveRefreshToken(ctx, userID); err != nil {
			return fmt.Errorf("refresh token removing error: %w", err)
		}

		if err := p.CreateRefreshToken(ctx, newtokenHash, userID); err != nil {
			return fmt.Errorf("refresh token creating error: %w", err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
