package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateChat updates chat with correct id.
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

deleteLoop:
	for _, id := range usersInChatIDs {
		for _, user := range chat.Users {
			if user.ID == id {
				continue deleteLoop
			}
		}

		if err := r.db.Exec("DELETE FROM user_chats WHERE chat_id = ? AND user_id = ?", chat.ID, id).Error; err != nil {
			return fmt.Errorf("user removing from chat error: %w", err)
		}
	}

addLoop:
	for _, user := range chat.Users {
		for _, id := range usersInChatIDs {
			if id == user.ID {
				continue addLoop
			}
		}

		if err := r.db.Exec("INSERT INTO user_chats (chat_id, user_id) VALUES (?, ?)", chat.ID, user.ID).Error; err != nil {
			return fmt.Errorf("user adding to chat error: %w", err)
		}
	}

	return nil
}

func (r *PostgresRepo) updateChatMessages(chat *domain.Chat) error {
	messagesInChat := make([]domain.Message, 0)

	if err := r.db.Raw("SELECT * FROM messages WHERE chat_id = ?", chat.ID).Scan(&messagesInChat).Error; err != nil {
		return fmt.Errorf("chat messages selection error: %w", err)
	}

	//Потом добавить удаление сообщений
	/*deleteLoop:
	for _, messageInChat := range messagesInChat {
		for _, newMessage := range chat.Messages {
			if messageInChat == newMessage {
				continue deleteLoop
			}
		}

		if err := r.db.Exec("DELETE FROM messages WHERE id = ?", messageInChat.ID).Error; err != nil {
			return fmt.Errorf("message deleting error: %w", err)
		}
	}
	*/

addLoop:
	for _, newMessage := range chat.Messages {
		for _, messageInChat := range messagesInChat {
			if newMessage == messageInChat {
				continue addLoop
			}
		}

		if err := r.db.Exec("INSERT INTO messages (chat_id, user_id, text, created_at) VALUES (?, ?, ?, ?)", newMessage.ChatID, newMessage.UserID, newMessage.Text, newMessage.CreatedAt).Error; err != nil {
			return fmt.Errorf("message adding to chat error: %w", err)
		}
	}

	return nil
}
