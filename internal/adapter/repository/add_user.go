package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// AddUser adds new user (domain.User) into the users table.
func (r *PostgresRepo) AddUser(user domain.User) (*domain.User, error) {
	if err := r.db.Model(&domain.User{}).Create(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, fmt.Errorf("some of user fields are duplicated: %w", domain.ErrDuplicated)
		}
		return nil, err
	}

	return &user, nil
}
