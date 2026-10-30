package usecase

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// GetMessages implements the get_messages endpoint on the business layer.
// The method receives dto.GetMessagesInput and returns dto.GetMessagesOutput and error.
// It returns all messages from chat.
// If the uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
func (u *Usecase) GetMessages(ctx context.Context, input dto.GetMessagesInput) (dto.GetMessagesOutput, error) {
	if ok, err := u.isUserInChat(ctx, input.UserID, input.ChatID); !ok {
		if err != nil {
			return dto.GetMessagesOutput{}, fmt.Errorf("user check error: %w", err)
		}
		return dto.GetMessagesOutput{}, domain.ErrUserIsNotInChat
	}

	messages, err := u.postgres.GetMessages(ctx, input.ChatID)
	if err != nil {
		return dto.GetMessagesOutput{}, fmt.Errorf("messages getting error: %w", err)
	}

	return dto.GetMessagesOutput{Messages: messages}, nil
}
