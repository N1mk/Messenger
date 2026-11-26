package usecase

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// AddChatInvite implements the add_chat_invite endpoint on the busines layer.
// The method recievs dto.AddChatInviteInput and context.Context and returns dto.AddChatInviteOutput and error.
// It checks if user with id = userID in chat and adds invite to user with id = InvitedID.
// If user with id = userID is not in chat it returns domain.ErrUserIsNotInChat.
func (u *Usecase) AddChatInvite(ctx context.Context, input dto.AddChatInviteInput) (dto.AddChatInviteOutput, error) {
	if ok, err := u.isUserInChat(input.UserID, input.ChatID); !ok {
		if err != nil {
			return dto.AddChatInviteOutput{}, fmt.Errorf("inviter chat chack error: %w", err)
		}

		return dto.AddChatInviteOutput{}, fmt.Errorf("inviter is not in chat: %w", domain.ErrUserIsNotInChat)
	}

	if err := u.r.AddChatInvite(domain.ChatInvite{ChatID: input.ChatID, UserID: input.InvitedID}); err != nil {
		return dto.AddChatInviteOutput{}, fmt.Errorf("invite adding error: %w", err)
	}

	return dto.AddChatInviteOutput{}, nil
}
