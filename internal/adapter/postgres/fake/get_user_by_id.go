package postgres

import "messenger/internal/domain"

func (p *FakePostgres) GetUserByID(id int) (*domain.User, error) {
	return &p.Users[id], nil
}
