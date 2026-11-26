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
func (u *Usecase) CreateChat(ctx context.Context, input dto.CreateChatInput) (output dto.CreateChatOutput, err error) {
	u.r.Begin(ctx)
	defer func() {
		if err != nil {
			u.r.Rollback()
			return
		}
		u.r.Commit()
	}()

	creator, err := u.r.GetUserByID(input.UserID)
	if err != nil {
		u.l.LogError(fmt.Sprintf("creator getting error: %s", err.Error()))
		return dto.CreateChatOutput{}, fmt.Errorf("creator getting error: %w", err)
	}

	chat, err := u.r.AddChat(domain.Chat{Name: input.Name, OwnerID: input.UserID, Users: []domain.User{*creator}})
	if err != nil {
		u.l.LogError(fmt.Sprintf("chat creating error: %s", err.Error()))
		return dto.CreateChatOutput{}, fmt.Errorf("chat creating error: %w", err)
	}

	return dto.CreateChatOutput{ChatID: chat.ID}, nil
}
