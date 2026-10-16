package usecase

import (
	"fmt"
	"messenger/internal/dto"
)

// DeleteAccount implements the delete_account endpoint on the business layer.
// The method receives dto.DeleteAccountInput and returns dto.DeleteAccountOutput and error.
// It deletes user by id.
func (u *Usecase) DeleteAccount(input dto.DeleteAccountInput) (dto.DeleteAccountOutput, error) {
	user, err := u.postgres.GetUserByID(input.UserID)
	if err != nil {
		return dto.DeleteAccountOutput{}, fmt.Errorf("user getting error: %w", err)
	}

	for _, chat := range user.Chats {
		chatPointer, err := u.postgres.GetChatByID(chat.ID)
		if err != nil {
			return dto.DeleteAccountOutput{}, fmt.Errorf("chat getting error: %w", err)
		}

		chat = *chatPointer

		if len(chat.Users) == 1 {
			if _, err := u.DeleteChat(dto.DeleteChatInput{ChatID: chat.ID, UserID: input.UserID}); err != nil {
				return dto.DeleteAccountOutput{}, err
			}
		}
	}

	if err := u.postgres.RemoveUser(input.UserID); err != nil {
		return dto.DeleteAccountOutput{}, fmt.Errorf("user removing error: %w", err)
	}

	return dto.DeleteAccountOutput{}, nil
}
