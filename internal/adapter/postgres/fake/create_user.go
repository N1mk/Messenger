package postgres

import (
	"context"
	"messenger/internal/domain"
)

func (p *FakePostgres) CreateUser(_ context.Context, user domain.User) (*domain.User, error) {
	user.ID = len(p.Users)

	p.Users = append(p.Users, user)

	return &user, nil
}
