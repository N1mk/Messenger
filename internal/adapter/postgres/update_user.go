package postgres

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// UpdateUser updates user with correct id.
func (p *Postgres) UpdateUser(user *domain.User) error {
	var userWithChats domain.User
	if err := p.db.Preload("Chats").First(&userWithChats, user.ID).Error; err != nil {
		return fmt.Errorf("user getting error: %w", err)
	}

	for _, chat := range userWithChats.Chats {
		var chatWithUsers domain.Chat
		if err := p.db.Preload("Users").First(&chatWithUsers, chat.ID).Error; err != nil {
			return fmt.Errorf("chat without users searching error: %w", err)
		}

		if len(chatWithUsers.Users) <= 1 {
			if err := p.RemoveChat(chatWithUsers.ID); err != nil {
				return fmt.Errorf("chat without users deleting error: %w", err)
			}
		}
	}

	if err := p.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user updating error: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("user updating error: %w", err)
	}

	return nil
}
