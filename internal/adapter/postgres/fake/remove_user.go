package postgres

func (p *FakePostgres) RemoveUser(id int) error {
	p.Users = append(p.Users[:id], p.Users[id+1:]...)

	return nil
}
