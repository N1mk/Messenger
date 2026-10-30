package postgres

import (
	"context"
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetUserByID selects user from the users table by id.
func (p *Postgres) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	err := p.db.WithContext(ctx).Preload("Chats").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
