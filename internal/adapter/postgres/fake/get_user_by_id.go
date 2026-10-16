package postgres

import "messenger/internal/domain"

func (p *FakePostgres) GetUserByID(id uint) (*domain.User, error) {
	return &p.Users[id], nil
}
