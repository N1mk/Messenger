package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// AddUserToChat adds user with id = userID to chat with id = chatID. If user with id = userID or chat with id = chatID not found AddUserToChat returns domain.ErrNotFound.
func (pr *PostgresRepo) AddUserToChat(userID int, chatID int) error {
	if err := pr.db.Raw("SELECT * FROM users WHERE id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("user exist check error: %w", err)
	}

	if err := pr.db.Raw("SELECT * FROM chats WHERE id = ?", chatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("chat not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("chat exist check error: %w", err)
	}

	if err := pr.db.Exec("INSERT INTO user_chats (chat_id, user_id) VALUES (?, ?)", chatID, userID).Error; err != nil {
		return fmt.Errorf("user adding error: %w", err)
	}

	return nil
}
