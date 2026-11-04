package repository

import (
	"errors"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveChat removes chat from the chats table by id.
func (r *PostgresRepo) RemoveChat(id int) error {
	if err := r.db.Exec("DELETE FROM user_chats WHERE chat_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	if err := r.db.Exec("DELETE FROM messages WHERE chat_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	if err := r.db.Exec("DELETE FROM chats WHERE id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	return nil
}
