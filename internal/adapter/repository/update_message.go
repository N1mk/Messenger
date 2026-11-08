package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateMessage updates all non-related fields of message (domain.Message) with id = message.ID. If message with id = message.ID not found UpdateMessage returns domain.ErrNotFound.
func (r *PostgresRepo) UpdateMessage(message *domain.Message) error {
	if err := r.db.Model(&domain.Message{}).Where("id = ?", message.ID).Updates(&domain.Message{Text: message.Text}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("message not found: %w", domain.ErrNotFound)
		}
		return err
	}

	return nil
}
