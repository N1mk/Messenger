package repository

import (
	"fmt"
	"messenger/internal/domain"
)

// СreateChat creates new element in the chats table.
func (r *PostgresRepo) CreateChat(chat domain.Chat) (*domain.Chat, error) {
	var chatID int

	if err := r.db.Raw("INSERT INTO chats (owner_id, name) VALUES (?, ?) RETURNING id", chat.OwnerID, chat.Name).Scan(&chatID).Error; err != nil {
		return nil, fmt.Errorf("chat creating error: %w", err)
	}

	chat.ID = chatID

	for _, user := range chat.Users {
		if err := r.db.Exec("INSERT INTO user_chats (chat_id, user_id) VALUES (?, ?)", chat.ID, user.ID).Error; err != nil {
			return nil, fmt.Errorf("user adding to chat error: %w", err)
		}
	}

	return &chat, nil
}
