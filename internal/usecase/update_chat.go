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
	if ok, _ := u.isUserInChat(input.UserID, input.ChatID); !ok {
		return dto.UpdateChatOutput{}, domain.ErrUnauthorized
	}

	if err := u.postgres.UpdateChat(&domain.Chat{ID: input.UserID, Name: input.Name}); err != nil {
		return dto.UpdateChatOutput{}, fmt.Errorf("chat updating error: %w", err)
	}

	return dto.UpdateChatOutput{}, nil
}
