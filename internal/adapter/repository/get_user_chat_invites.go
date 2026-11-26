package repository

import (
	"errors"
	"fmt"
	"messenger/internal/domain"

	"gorm.io/gorm"
)

// GetUserChatInvites returns slice of all user's chat invites. If user with id = userID not found GetUserChatInvites returns domain.ErrNotFound.
func (pr *PostgresRepo) GetUserChatInvites(userID int) (invites []*domain.ChatInvite, err error) {
	if err := pr.db.Raw("SELECT * FROM users WHERE id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*domain.ChatInvite{}, fmt.Errorf("user not found: %w", domain.ErrNotFound)
		}
		return []*domain.ChatInvite{}, fmt.Errorf("user exist check error: %w", err)
	}

	if err := pr.db.Raw("SELECT * FROM chat_invites WHERE user_id = ?", userID).Scan(&invites).Error; err != nil {
		return []*domain.ChatInvite{}, fmt.Errorf("chat invites getting error: %w", err)
	}

	return invites, nil
}
