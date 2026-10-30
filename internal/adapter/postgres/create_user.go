package postgres

import (
	"context"
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// СreateUser creates new element in the users table.
func (p *Postgres) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	err := p.db.WithContext(ctx).Model(&domain.User{}).Create(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain.ErrDuplicated
		}
		return nil, err
	}

	return &user, nil
}
