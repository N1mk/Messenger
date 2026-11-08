package usecase

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"time"
)

// SendMessage implements the send_message endpoint on the business layer.
// The method receives dto.SendMessageInput and returns dto.SendMessageOutput and error.
// It creates new message in chat with ID = input.ChatID.
// If the uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
func (u *Usecase) SendMessage(ctx context.Context, input dto.SendMessageInput) (output dto.SendMessageOutput, err error) {
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
			return dto.SendMessageOutput{}, fmt.Errorf("user membership check error: %w", err)
		}
		u.l.LogError(fmt.Sprintf("user in not in the chat: %s", domain.ErrUserIsNotInChat.Error()))
		return dto.SendMessageOutput{}, domain.ErrUserIsNotInChat
	}

	message, err := u.r.AddMessage(&domain.Message{ChatID: input.ChatID, UserID: input.UserID, Text: input.Text, CreatedAt: time.Now()})
	if err != nil {
		return dto.SendMessageOutput{}, fmt.Errorf("message creating error: %w", &err)
	}

	return dto.SendMessageOutput{ID: message.ID}, nil
}
