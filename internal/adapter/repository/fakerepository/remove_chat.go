package fakerepository

func (r *FakePostgresRepo) RemoveChat(id int) error {
	r.Chats = append(r.Chats[:id], r.Chats[id+1:]...)

	return nil
}
