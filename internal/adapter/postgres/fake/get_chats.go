package postgres

func (p *FakePostgres) GetChatsCount() int {
	return len(p.Chats)
}
