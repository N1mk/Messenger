package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateUser updates user with correct id.
func (r *PostgresRepo) UpdateUser(user *domain.User) error {
	if err := r.db.Updates(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user updating error: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("user updating error: %w", err)
	}

	return nil
}
