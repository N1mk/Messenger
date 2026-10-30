package postgres

import "context"

func (p *FakePostgres) GetChatsCount(_ context.Context) int {
	return len(p.Chats)
}
