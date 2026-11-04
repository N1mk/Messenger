package fakerepository

import (
	"messenger/internal/domain"
)

func (r *FakePostgresRepo) CreateChat(chat domain.Chat) (*domain.Chat, error) {
	chat.ID = len(r.Chats)

	r.Chats = append(r.Chats, chat)

	return &chat, nil
}
