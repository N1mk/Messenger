package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetUserByName recives user's username and returns *domain.User with same username or domain.ErrNotFound if user is not found.
func (r *PostgresRepo) GetUserByName(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Chats").Raw("SELECT * FROM users WHERE username = ?", username).Scan(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return nil, err
	}

	return &user, nil
}
