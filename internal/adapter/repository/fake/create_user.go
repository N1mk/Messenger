package postgres

import (
	"messenger/internal/domain"
)

func (r *FakePostgresRepo) CreateUser(user domain.User) (*domain.User, error) {
	user.ID = len(r.Users)

	r.Users = append(r.Users, user)

	return &user, nil
}
