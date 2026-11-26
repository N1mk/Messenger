package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateChat updates all non-related fields of chat (domain.Chat) with id = chat.ID. If chat with id = chat.ID not found UpdateChat returns domain.ErrNotFound.
func (pr *PostgresRepo) UpdateChat(chat domain.Chat) error {
	if err := pr.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&domain.Chat{}).Where("id = ?", chat.ID).Updates(&domain.Chat{OwnerID: chat.OwnerID, Name: chat.Name}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("chat basic fields updating error: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("chat basic fields updating error: %w", err)
	}

	return nil
}
