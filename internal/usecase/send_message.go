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
func (u *Usecase) SendMessage(ctx context.Context, input dto.SendMessageInput) (dto.SendMessageOutput, error) {
	if ok, err := u.isUserInChat(ctx, input.UserID, input.ChatID); !ok {
		if err != nil {
			return dto.SendMessageOutput{}, fmt.Errorf("user check error: %w", err)
		}
		return dto.SendMessageOutput{}, domain.ErrUserIsNotInChat
	}

	chat, err := u.postgres.GetChatByID(ctx, input.ChatID)
	if err != nil {
		return dto.SendMessageOutput{}, fmt.Errorf("chat getting error: %w", err)
	}

	chat.Messages = append(chat.Messages, domain.Message{ChatID: chat.ID, UserID: input.UserID, Text: input.Text, CreatedAt: time.Now()})

	if err := u.postgres.UpdateChat(ctx, chat); err != nil {
		return dto.SendMessageOutput{}, fmt.Errorf("chat update error: %w", err)
	}

	return dto.SendMessageOutput{}, nil
}
