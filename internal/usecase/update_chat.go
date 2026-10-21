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

	if chat.Owner.ID != input.UserID {
		return dto.UpdateChatOutput{}, domain.ErrUnauthorized
	}

	owner, err := u.postgres.GetUserByID(input.OwnerID)
	if err != nil {
		return dto.UpdateChatOutput{}, fmt.Errorf("new chat owner getting error: %w", err)
	}

	if err := u.postgres.UpdateChat(&domain.Chat{ID: input.UserID, Owner: *owner, Name: input.Name}); err != nil {
		return dto.UpdateChatOutput{}, fmt.Errorf("chat updating error: %w", err)
	}

	return dto.UpdateChatOutput{}, nil
}
