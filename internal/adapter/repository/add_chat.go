package repository

import (
	"fmt"
	"messenger/internal/domain"
)

// AddChat adds new chat (domain.Chat) into the chats table.
func (r *PostgresRepo) AddChat(chat domain.Chat) (*domain.Chat, error) {
	var chatID int

	if err := r.db.Raw("INSERT INTO chats (owner_id, name) VALUES (?, ?) RETURNING id", chat.OwnerID, chat.Name).Scan(&chatID).Error; err != nil {
		return nil, err
	}

	chat.ID = chatID

	var placeholders string
	values := make([]interface{}, 0, len(chat.Users)*2)

	for _, user := range chat.Users {
		if placeholders != "" {
			placeholders += ", "
		}

		placeholders += "(?, ?)"

		values = append(values, chat.ID, user.ID)
	}

	if err := r.db.Exec(fmt.Sprintf("INSERT INTO user_chats (chat_id, user_id) VALUES %s", placeholders), values...).Error; err != nil {
		return nil, fmt.Errorf("user adding to chat error: %w", err)
	}

	return &chat, nil
}
