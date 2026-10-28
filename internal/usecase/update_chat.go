package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// UpdateChat implements the update_chat endpoint on the business layer.
// The method receives dto.UpdateChatInput and returns dto.UpdateChatOutput and error.
// It updates chat with same id as in input.
func (u *Usecase) UpdateChat(input dto.UpdateChatInput) (dto.UpdateChatOutput, error) {
	chat, err := u.postgres.GetChatByID(input.ChatID)
	if err != nil {
		return dto.UpdateChatOutput{}, fmt.Errorf("chat getting error: %w", err)
	}

	if input.UserID != chat.OwnerID {
		return dto.UpdateChatOutput{}, domain.ErrUnauthorized
	}

	users := make([]domain.User, len(input.Usernames))
	for i, username := range input.Usernames {
		user, err := u.postgres.GetUserByName(username)
		if err != nil {
			return dto.UpdateChatOutput{}, fmt.Errorf("user getting error: %w", err)
		}

		users[i] = *user
	}

	if err := u.postgres.UpdateChat(&domain.Chat{ID: input.ChatID, OwnerID: input.OwnerID, Name: input.Name, Users: users}); err != nil {
		return dto.UpdateChatOutput{}, fmt.Errorf("chat updating error: %w", err)
	}

	return dto.UpdateChatOutput{}, nil
}
