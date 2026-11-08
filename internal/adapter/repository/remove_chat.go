package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveChat removes chat (domain.Chat) with given id from the chats table. If chat with given id not found RemoveChat returns domain.ErrNotFound.
func (r *PostgresRepo) RemoveChat(id int) error {
	if err := r.db.Exec("DELETE FROM user_chats WHERE chat_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user chat relations remove error: chat not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("user chat relations remove error: %w", err)
	}

	if err := r.db.Exec("DELETE FROM messages WHERE chat_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("chat messages remove error: chat not found: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("chat messages remove error: %w", err)
	}

	if err := r.db.Exec("DELETE FROM chats WHERE id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("chat not found: %w", domain.ErrNotFound)
		}
		return err
	}

	return nil
}
