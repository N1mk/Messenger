package postgres

import (
	"context"
	"messenger/internal/domain"
)

func (p *FakePostgres) GetUserByName(_ context.Context, username string) (*domain.User, error) {
	for _, user := range p.Users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, domain.ErrNotFound
}
