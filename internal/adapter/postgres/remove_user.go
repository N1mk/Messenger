package postgres

import (
	"context"
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveUser removes user from the users table by id.
func (p *Postgres) RemoveUser(ctx context.Context, id int) error {
	var chatsToRemove []domain.Chat

	var user domain.User
	if err := p.db.WithContext(ctx).Preload("Chats").First(&user, id).Error; err != nil {
		return fmt.Errorf("user getting error: %w", err)
	}

	for _, chat := range user.Chats {
		var chatWithUsers domain.Chat
		if err := p.db.WithContext(ctx).Preload("Users").First(&chatWithUsers, chat.ID).Error; err != nil {
			return fmt.Errorf("chat without users searching error: %w", err)
		}

		if len(chatWithUsers.Users) <= 1 {
			chatsToRemove = append(chatsToRemove, chatWithUsers)
		}
	}

	if err := p.db.WithContext(ctx).Exec("DELETE FROM user_chats WHERE user_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	if err := p.db.WithContext(ctx).Delete(&domain.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	for _, chat := range chatsToRemove {
		if err := p.RemoveChat(ctx, chat.ID); err != nil {
			return fmt.Errorf("chat without users deleting error: %w", err)
		}
	}

	return nil
}
