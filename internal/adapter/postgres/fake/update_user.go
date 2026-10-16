package postgres

import "messenger/internal/domain"

func (p *FakePostgres) UpdateUser(user *domain.User) error {
	p.Users[user.ID] = *user

	return nil
}
