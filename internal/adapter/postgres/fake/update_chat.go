package postgres

import "messenger/internal/domain"

func (p *FakePostgres) UpdateChat(chat *domain.Chat) error {
	p.Chats[chat.ID] = *chat

	return nil
}
