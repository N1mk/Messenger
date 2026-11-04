package fakerepository

import (
	"context"
	"messenger/internal/domain"
)

type FakePostgresRepo struct {
	Users []domain.User
	Chats []domain.Chat
}

func NewFakePostgres() *FakePostgresRepo {
	users := make([]domain.User, 0)
	chats := make([]domain.Chat, 0)
	return &FakePostgresRepo{Users: users, Chats: chats}
}

// Заглушка
func (r *FakePostgresRepo) Begin(ctx context.Context) {}

// Заглушка
func (r *FakePostgresRepo) Rollback() {}

// Заглушка
func (r *FakePostgresRepo) Commit() {}
