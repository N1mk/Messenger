package postgres

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// СreateUser creates new element in the users table.
func (p *Postgres) CreateUser(user domain.User) (*domain.User, error) {
	err := p.db.Model(&domain.User{}).Create(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain.ErrDuplicated
		}
		return nil, err
	}

	return &user, nil
}
