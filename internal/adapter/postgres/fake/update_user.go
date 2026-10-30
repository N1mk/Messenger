package postgres

import (
	"context"
	"messenger/internal/domain"
)

func (p *FakePostgres) UpdateUser(_ context.Context, user *domain.User) error {
	p.Users[user.ID] = *user

	return nil
}
