package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"
	"strconv"

	"gorm.io/gorm"
)

// UpdateChat updates all non-related fields of chat (domain.Chat) with id = chat.ID. If chat with id = chat.ID not found UpdateChat returns domain.ErrNotFound.
func (r *PostgresRepo) UpdateChat(chat *domain.Chat) error {
	if err := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&domain.Chat{}).Where("id = ?", chat.ID).Updates(&domain.Chat{OwnerID: chat.OwnerID, Name: chat.Name}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("chat basic fields updating error: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("chat basic fields updating error: %w", err)
	}

	if err := r.updateChatUsers(chat); err != nil {
		return fmt.Errorf("chat users updating error: %w", err)
	}

	if err := r.updateChatMessages(chat); err != nil {
		return fmt.Errorf("chat messages updating error: %w", err)
	}

	return nil
}

func (r *PostgresRepo) updateChatUsers(chat *domain.Chat) error {
	usersInChatIDs := make([]int, 0)

	if err := r.db.Raw("SELECT user_id FROM user_chats WHERE chat_id = ?", chat.ID).Scan(&usersInChatIDs).Error; err != nil {
		return fmt.Errorf("chat users selection error: %w", err)
	}

	var chatIDsDelete string
	var userIDsDelete string

deleteLoop:
	for _, id := range usersInChatIDs {
		for _, user := range chat.Users {
			if user.ID == id {
				continue deleteLoop
			}
		}

		if chatIDsDelete != "" {
			chatIDsDelete += ", "
		}
		chatIDsDelete += strconv.Itoa(chat.ID)

		if userIDsDelete != "" {
			userIDsDelete += ", "
		}
		userIDsDelete += strconv.Itoa(id)
	}

	if chatIDsDelete != "" && userIDsDelete != "" {
		if err := r.db.Exec(fmt.Sprintf("DELETE FROM user_chats WHERE chat_id IN (%s) AND user_id IN (%s)", chatIDsDelete, userIDsDelete)).Error; err != nil {
			return fmt.Errorf("users deleting from chat error: %w", err)
		}
	}

	var userChatsAdd string

addLoop:
	for _, user := range chat.Users {
		for _, id := range usersInChatIDs {
			if id == user.ID {
				continue addLoop
			}
		}

		if userChatsAdd != "" {
			userChatsAdd += ", "
		}
		userChatsAdd += fmt.Sprintf("(%s, %s)", strconv.Itoa(chat.ID), strconv.Itoa(user.ID))
	}

	if userChatsAdd != "" {
		if err := r.db.Exec(fmt.Sprintf("INSERT INTO user_chats (chat_id, user_id) VALUES %s", userChatsAdd)).Error; err != nil {
			return fmt.Errorf("users adding to chat error: %w", err)
		}
	}

	return nil
}

func (r *PostgresRepo) updateChatMessages(chat *domain.Chat) error {
	messagesInChat := make([]domain.Message, 0)

	if err := r.db.Raw("SELECT * FROM messages WHERE chat_id = ?", chat.ID).Scan(&messagesInChat).Error; err != nil {
		return fmt.Errorf("chat messages selection error: %w", err)
	}

	var placeholders string
	values := make([]interface{}, 0)

addLoop:
	for _, newMessage := range chat.Messages {
		for _, messageInChat := range messagesInChat {
			if newMessage == messageInChat {
				continue addLoop
			}
		}

		if placeholders != "" {
			placeholders += ", "
		}

		placeholders += "(?, ?, ?, ?)"
		values = append(values, newMessage.ChatID, newMessage.UserID, newMessage.Text, newMessage.CreatedAt)
	}

	if placeholders != "" {
		if err := r.db.Exec(fmt.Sprintf("INSERT INTO messages (chat_id, user_id, text, created_at) VALUES %s", placeholders), values...).Error; err != nil {
			return fmt.Errorf("messages adding to chat error: %w", err)
		}
	}

	return nil
}
