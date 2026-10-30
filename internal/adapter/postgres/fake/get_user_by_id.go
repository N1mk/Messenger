package postgres

import (
	"context"
	"messenger/internal/domain"
)

func (p *FakePostgres) GetUserByID(_ context.Context, id int) (*domain.User, error) {
	return &p.Users[id], nil
}
