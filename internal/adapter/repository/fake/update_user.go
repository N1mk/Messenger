package postgres

import (
	"messenger/internal/domain"
)

func (r *FakePostgresRepo) UpdateUser(user *domain.User) error {
	r.Users[user.ID] = *user

	return nil
}
