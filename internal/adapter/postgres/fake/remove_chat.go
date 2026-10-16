package postgres

func (p *FakePostgres) RemoveChat(id uint) error {
	p.Chats = append(p.Chats[:id], p.Chats[id+1:]...)

	return nil
}
