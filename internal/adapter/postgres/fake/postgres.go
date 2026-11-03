package postgres

import (
	"context"
	"messenger/internal/domain"
)

type FakePostgres struct {
	Users []domain.User
	Chats []domain.Chat
}

func NewFakePostgres() *FakePostgres {
	users := make([]domain.User, 0)
	chats := make([]domain.Chat, 0)
	return &FakePostgres{Users: users, Chats: chats}
}

// Заглушка
func (p *FakePostgres) Begin(ctx context.Context) {}

// Заглушка
func (p *FakePostgres) Rollback() {}

// Заглушка
func (p *FakePostgres) Commit() {}
