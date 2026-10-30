package postgres

import (
	"context"
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveRefreshToken removes refresh token from refresh_tokens table.
func (p *Postgres) RemoveRefreshToken(ctx context.Context, userID int) error {
	if err := p.db.WithContext(ctx).Exec("DELETE FROM refresh_tokens WHERE user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
