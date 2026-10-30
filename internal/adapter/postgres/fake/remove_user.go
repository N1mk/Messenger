package postgres

import "context"

func (p *FakePostgres) RemoveUser(_ context.Context, id int) error {
	p.Users = append(p.Users[:id], p.Users[id+1:]...)

	return nil
}
