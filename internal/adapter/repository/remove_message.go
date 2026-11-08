package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveMessage removes message (domain.Message) with given id from the messages table. If message with given id not found RemoveMessage returns domain.ErrNotFound.
func (r *PostgresRepo) RemoveMessage(chatID int) error {
	if err := r.db.Exec("DELETE FROM messages WHERE chat_id = ?", chatID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("chat not found: %w", domain.ErrNotFound)
		}
		return err
	}

	return nil
}
