package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// RemoveUser removes user (domain.User) with given id from the users table. If user with given id not found RemoveUser returns domain.ErrNotFound.
func (pr *PostgresRepo) RemoveUser(id int) error {
	var user domain.User
	if err := pr.db.Preload("Chats.Users").First(&user, id).Error; err != nil {
		return fmt.Errorf("user getting error: %w", err)
	}

	var chatsToRemoveIDs []int

	for _, chat := range user.Chats {
		if len(chat.Users) <= 1 {
			chatsToRemoveIDs = append(chatsToRemoveIDs, chat.ID)
			continue
		}

		if chat.OwnerID == id {
			var newOwner domain.User
			for _, member := range chat.Users {
				if member.ID != user.ID {
					newOwner = member
					break
				}
			}

			chat.OwnerID = newOwner.ID

			if err := pr.UpdateChat(chat); err != nil {
				return fmt.Errorf("chat owner changing error: %w", err)
			}
		}
	}

	if err := pr.db.Exec("DELETE FROM user_chats WHERE user_id = ?", id).Error; err != nil {
		return fmt.Errorf("chats constraint deleting error: %w", err)
	}

	if err := pr.db.Exec("DELETE FROM refresh_tokens WHERE user_id = ?", id).Error; err != nil {
		return fmt.Errorf("refresh tokens constraint deleting error: %w", err)
	}

	if err := pr.db.Delete(&domain.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	if len(chatsToRemoveIDs) != 0 {
		if err := pr.db.Delete(&domain.Chat{}, chatsToRemoveIDs).Error; err != nil {
			return fmt.Errorf("empty chats removing error: %w", err)
		}
	}

	return nil
}
