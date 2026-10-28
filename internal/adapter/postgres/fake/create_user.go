package postgres

import "messenger/internal/domain"

func (p *FakePostgres) CreateUser(user domain.User) (*domain.User, error) {
	user.ID = len(p.Users)

	p.Users = append(p.Users, user)

	return &user, nil
}
