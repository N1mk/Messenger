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
	u.p.Begin(ctx)
	defer func() {
		if err != nil {
			u.p.Rollback()
			return
		}
		u.p.Commit()
	}()

	creator, err := u.p.GetUserByID(input.UserID)
	if err != nil {
		return dto.CreateChatOutput{}, fmt.Errorf("user getting error: %w", err)
	}

	users := make([]domain.User, len(input.Usernames)+1)
	users[0] = *creator

	for i, username := range input.Usernames {
		user, err := u.p.GetUserByName(username)
		if err != nil {
			return dto.CreateChatOutput{}, fmt.Errorf("user getting error: %w", err)
		}

		users[i+1] = *user
	}

	chat, err := u.p.CreateChat(domain.Chat{Name: input.Name, OwnerID: input.UserID, Users: users})
	if err != nil {
		return dto.CreateChatOutput{}, fmt.Errorf("chat creation error: %w", err)
	}

	return dto.CreateChatOutput{ChatID: chat.ID}, nil
}
