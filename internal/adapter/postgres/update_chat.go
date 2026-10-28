package postgres

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateChat updates chat with correct id.
func (p *Postgres) UpdateChat(chat *domain.Chat) error {
	if err := p.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&domain.Chat{}).Where("id = ?", chat.ID).Updates(&domain.Chat{ID: chat.ID, OwnerID: chat.OwnerID, Name: chat.Name}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("chat basic fields updating error: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("chat basic fields updating error: %w", err)
	}

	userIDsInChat := make([]int, 0)

	if err := p.db.Raw("SELECT user_id FROM user_chats WHERE chat_id = ?", chat.ID).Scan(&userIDsInChat).Error; err != nil {
		return fmt.Errorf("chat users selection error: %w", err)
	}

deleteLoop:
	for _, id := range userIDsInChat {
		for _, user := range chat.Users {
			if int(user.ID) == id {
				continue deleteLoop
			}
		}
		if err := p.db.Exec("DELETE FROM user_chats WHERE chat_id = ? AND user_id = ?", chat.ID, id).Error; err != nil {
			return fmt.Errorf("user removing from chat error: %w", err)
		}
	}

addLoop:
	for _, user := range chat.Users {
		for _, id := range userIDsInChat {
			if id == int(user.ID) {
				continue addLoop
			}
		}
		if err := p.db.Exec("INSERT INTO user_chats (chat_id, user_id) VALUES (?, ?)", chat.ID, user.ID).Error; err != nil {
			return fmt.Errorf("user adding to chat error: %w", err)
		}
	}

	return nil
}
