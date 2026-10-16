package postgres

import "messenger/internal/domain"

type FakePostgres struct {
	Users []domain.User
	Chats []domain.Chat
}

func NewFakePostgres() *FakePostgres {
	users := make([]domain.User, 0)
	chats := make([]domain.Chat, 0)
	return &FakePostgres{Users: users, Chats: chats}
}
