package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// DeleteChat implements the delete_chat endpoint on the business layer.
// The method receives dto.DeleteChatInput and returns dto.DeleteChatOutput and error.
// It deletes chat by id.
// If uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
func (u *Usecase) DeleteChat(input dto.DeleteChatInput) (dto.DeleteChatOutput, error) {
	if ok, err := u.isUserInChat(input.UserID, input.ChatID); !ok {
		if err != nil {
			return dto.DeleteChatOutput{}, fmt.Errorf("user check error: %w", err)
		}
		return dto.DeleteChatOutput{}, domain.ErrUserIsNotInChat
	}

	if err := u.postgres.RemoveChat(input.ChatID); err != nil {
		return dto.DeleteChatOutput{}, fmt.Errorf("chat removing error: %w", err)
	}

	return dto.DeleteChatOutput{}, nil
}
