package postgres

import (
	"messenger/internal/domain"
)

func (r *FakePostgresRepo) GetUserByID(id int) (*domain.User, error) {
	return &r.Users[id], nil
}
