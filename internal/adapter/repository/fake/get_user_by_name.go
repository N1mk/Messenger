package postgres

import (
	"messenger/internal/domain"
)

func (r *FakePostgresRepo) GetUserByName(username string) (*domain.User, error) {
	for _, user := range r.Users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, domain.ErrNotFound
}
