package fakerepository

import (
	"messenger/internal/domain"
)

func (r *FakePostgresRepo) GetChatByID(id int) (*domain.Chat, error) {
	return &r.Chats[id], nil
}
