package postgres

func (p *FakePostgres) RemoveUser(id uint) error {
	p.Users = append(p.Users[:id], p.Users[id+1:]...)

	return nil
}
