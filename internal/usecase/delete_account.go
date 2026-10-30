package usecase

import (
	"context"
	"fmt"
	"messenger/internal/dto"
)

// DeleteAccount implements the delete_account endpoint on the business layer.
// The method receives dto.DeleteAccountInput and returns dto.DeleteAccountOutput and error.
// It deletes user by id.
func (u *Usecase) DeleteAccount(ctx context.Context, input dto.DeleteAccountInput) (dto.DeleteAccountOutput, error) {
	user, err := u.postgres.GetUserByID(ctx, input.UserID)
	if err != nil {
		return dto.DeleteAccountOutput{}, fmt.Errorf("user getting error: %w", err)
	}

	for _, chat := range user.Chats {
		chatPointer, err := u.postgres.GetChatByID(ctx, chat.ID)
		if err != nil {
			return dto.DeleteAccountOutput{}, fmt.Errorf("chat getting error: %w", err)
		}

		chat = *chatPointer

		if len(chat.Users) == 1 {
			if _, err := u.DeleteChat(ctx, dto.DeleteChatInput{ChatID: chat.ID, UserID: input.UserID}); err != nil {
				return dto.DeleteAccountOutput{}, err
			}
		}
	}

	if err := u.postgres.RemoveUser(ctx, input.UserID); err != nil {
		return dto.DeleteAccountOutput{}, fmt.Errorf("user removing error: %w", err)
	}

	return dto.DeleteAccountOutput{}, nil
}
