package postgres

import "messenger/internal/domain"

func (p *FakePostgres) GetUserByName(username string) (*domain.User, error) {
	for _, user := range p.Users {
		if user.Username == username {
			return &user, nil
		}
	}

	return nil, domain.ErrNotFound
}
