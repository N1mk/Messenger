package postgres

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetUserByName selects user from the users table by name.
func (p *Postgres) GetUserByName(username string) (*domain.User, error) {
	var user domain.User
	err := p.db.Preload("Chats").Raw("SELECT * FROM users WHERE username = ?", username).Scan(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
