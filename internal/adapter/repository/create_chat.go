package repository

import (
	"fmt"
	"messenger/internal/domain"
	"strings"
)

// СreateChat creates new element in the chats table.
func (r *PostgresRepo) CreateChat(chat domain.Chat) (*domain.Chat, error) {
	var chatID int

	if err := r.db.Raw("INSERT INTO chats (owner_id, name) VALUES (?, ?) RETURNING id", chat.OwnerID, chat.Name).Scan(&chatID).Error; err != nil {
		return nil, fmt.Errorf("chat creating error: %w", err)
	}

	chat.ID = chatID

	// Adding users to chat
	query := "INSERT INTO user_chats (chat_id, user_id) VALUES "
	values := make([]interface{}, 0, len(chat.Users)*2)
	placeholders := make([]string, 0, len(chat.Users))

	for _, user := range chat.Users {
		placeholders = append(placeholders, "(?, ?)")
		values = append(values, chat.ID, user.ID)
	}

	query += strings.Join(placeholders, ", ")

	if err := r.db.Exec(query, values...).Error; err != nil {
		return nil, fmt.Errorf("user adding to chat error: %w", err)
	}

	return &chat, nil
}
