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
func (u *Usecase) GetMessages(ctx context.Context, input dto.GetMessagesInput) (output dto.GetMessagesOutput, err error) {
	u.r.Begin(ctx)
	defer func() {
		if err != nil {
			u.r.Rollback()
			return
		}
		u.r.Commit()
	}()

	if ok, err := u.isUserInChat(input.UserID, input.ChatID); !ok {
		if err != nil {
			u.l.LogError(fmt.Sprintf("user membership check error: %s", err.Error()))
			return dto.GetMessagesOutput{}, fmt.Errorf("user membership check error: %w", err)
		}
		return dto.GetMessagesOutput{}, domain.ErrUserIsNotInChat
	}

	messages, err := u.r.GetChatMessages(input.ChatID, input.LastMessageID, input.PageSize)
	if err != nil {
		u.l.LogError(fmt.Sprintf("message getting error: %s", err.Error()))
		return dto.GetMessagesOutput{}, fmt.Errorf("messages getting error: %w", err)
	}

	return dto.GetMessagesOutput{Messages: messages}, nil
}
