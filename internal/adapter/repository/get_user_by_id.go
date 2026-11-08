package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetUserByID recives user's id and returns *domain.User with same id or domain.ErrNotFound if user is not found.
func (r *PostgresRepo) GetUserByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Chats").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return nil, err
	}

	return &user, nil
}
