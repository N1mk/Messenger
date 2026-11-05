package usecase

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// DeleteChat implements the delete_chat endpoint on the business layer.
// The method receives dto.DeleteChatInput and returns dto.DeleteChatOutput and error.
// It deletes chat by id.
// If uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
func (u *Usecase) DeleteChat(ctx context.Context, input dto.DeleteChatInput) (output dto.DeleteChatOutput, err error) {
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
			return dto.DeleteChatOutput{}, fmt.Errorf("user membership check error: %w", err)
		}
		u.l.LogError(fmt.Sprintf("user in not in the chat: %s", domain.ErrUserIsNotInChat.Error()))
		return dto.DeleteChatOutput{}, domain.ErrUserIsNotInChat
	}

	if err := u.r.RemoveChat(input.ChatID); err != nil {
		u.l.LogError(fmt.Sprintf("chat removing error: %s", err.Error()))
		return dto.DeleteChatOutput{}, fmt.Errorf("chat removing error: %w", err)
	}

	return dto.DeleteChatOutput{}, nil
}
