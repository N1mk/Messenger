package postgres

import (
	"messenger/internal/domain"
)

func (r *FakePostgresRepo) UpdateChat(chat *domain.Chat) error {
	r.Chats[chat.ID] = *chat

	return nil
}
