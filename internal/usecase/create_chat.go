package usecase

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// CreateChat implements the create_chat endpoint on the business layer.
// The method receives dto.CreateChatInput and returns dto.CreateChatOutput and error.
// It creates new chat.
func (u *Usecase) CreateChat(ctx context.Context, input dto.CreateChatInput) (dto.CreateChatOutput, error) {
	creator, err := u.p.GetUserByID(ctx, input.UserID)
	if err != nil {
		return dto.CreateChatOutput{}, fmt.Errorf("user getting error: %w", err)
	}

	chat, err := u.p.CreateChat(ctx, domain.Chat{Name: input.Name, OwnerID: input.UserID})
	if err != nil {
		return dto.CreateChatOutput{}, fmt.Errorf("chat creation error: %w", err)
	}

	creator.Chats = []domain.Chat{*chat}

	err = u.p.UpdateUser(ctx, creator)
	if err != nil {
		return dto.CreateChatOutput{}, fmt.Errorf("user update error: %w", err)
	}

	return dto.CreateChatOutput{ChatID: chat.ID}, nil
}
