package postgres

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateUser updates user with correct id.
func (p *Postgres) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := p.db.WithContext(ctx).Updates(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user updating error: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("user updating error: %w", err)
	}

	return nil
}
