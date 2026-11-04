package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveUser removes user from the users table by id.
func (r *PostgresRepo) RemoveUser(id int) error {
	var chatsToRemove []domain.Chat

	var user domain.User
	if err := r.db.Preload("Chats").First(&user, id).Error; err != nil {
		return fmt.Errorf("user getting error: %w", err)
	}

	for _, chat := range user.Chats {
		var chatWithUsers domain.Chat
		if err := r.db.Preload("Users").First(&chatWithUsers, chat.ID).Error; err != nil {
			return fmt.Errorf("chat without users searching error: %w", err)
		}

		if len(chatWithUsers.Users) <= 1 {
			chatsToRemove = append(chatsToRemove, chatWithUsers)
			continue
		}
		if chatWithUsers.OwnerID == id {
			var newOwner domain.User
			for _, member := range chatWithUsers.Users {
				if member.ID != user.ID {
					newOwner = member
					break
				}
			}

			chatWithUsers.OwnerID = newOwner.ID

			if err := r.UpdateChat(&chatWithUsers); err != nil {
				return fmt.Errorf("chat owner changing error: %w", err)
			}
		}
	}

	if err := r.db.Exec("DELETE FROM user_chats WHERE user_id = ?", id).Error; err != nil {
		return fmt.Errorf("chats constraint deleting error: %w", err)
	}

	if err := r.db.Exec("DELETE FROM refresh_tokens WHERE user_id = ?", id).Error; err != nil {
		return fmt.Errorf("refresh tokens constraint deleting error: %w", err)
	}

	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	for _, chat := range chatsToRemove {
		if err := r.RemoveChat(chat.ID); err != nil {
			return fmt.Errorf("chat without users deleting error: %w", err)
		}
	}

	return nil
}
