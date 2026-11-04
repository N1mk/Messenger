package fakerepository

func (r *FakePostgresRepo) RemoveUser(id int) error {
	r.Users = append(r.Users[:id], r.Users[id+1:]...)

	return nil
}
