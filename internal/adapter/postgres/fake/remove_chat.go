package postgres

import "context"

func (p *FakePostgres) RemoveChat(_ context.Context, id int) error {
	p.Chats = append(p.Chats[:id], p.Chats[id+1:]...)

	return nil
}
