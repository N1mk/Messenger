package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateUser updates all non-related fields of user (domain.User) with id = user.ID. If user with id = user.ID not found UpdateUser returns domain.ErrNotFound.
func (r *PostgresRepo) UpdateUser(user domain.User) error {
	if err := r.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(&domain.User{Username: user.Username, Email: user.Email, PasswordHash: user.PasswordHash}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return err
	}

	return nil
}
