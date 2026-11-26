package usecase

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// AddUserToChat implements the add_user_to_chat endpoint on the business layer.
// The method receives dto.AddUserToChatInput and context.Context and returns dto.AddUserToChatOutput and error.
// It adds user to chat and makes him listen the chat messages.
func (u *Usecase) AddUserToChat(ctx context.Context, input dto.AddUserToChatInput) (dto.AddUserToChatOutput, error) {
	invites, err := u.r.GetUserChatInvites(input.UserID)
	if err != nil {
		return dto.AddUserToChatOutput{}, fmt.Errorf("user invites getting error: %w", err)
	}

	invited := false
	for _, invite := range invites {
		if invite.ChatID == input.ChatID {
			invited = true
			break
		}
	}

	if !invited {
		return dto.AddUserToChatOutput{}, fmt.Errorf("user id not invited to chat: %w", domain.ErrUserIsNotInvited)
	}

	if err := u.r.AddUserToChat(input.UserID, input.ChatID); err != nil {
		return dto.AddUserToChatOutput{}, fmt.Errorf("adding user to chat error: %w", err)
	}

	return dto.AddUserToChatOutput{}, nil
}
