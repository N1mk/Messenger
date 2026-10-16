package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// GetChats implements the get_chats endpoint on the business layer.
// The method receives dto.GetChatsInput and returns dto.GetChatsOutput and error.
// It returns all chats, in which the user is a member.
func (u *Usecase) GetChats(input dto.GetChatsInput) (dto.GetChatsOutput, error) {
	user, err := u.postgres.GetUserByID(input.UserID)
	if err != nil {
		return dto.GetChatsOutput{}, fmt.Errorf("user getting error: %w", err)
	}

	chats := make([]*domain.Chat, 0)
	for _, chat := range user.Chats {
		chats = append(chats, &chat)
	}

	return dto.GetChatsOutput{Chats: chats}, nil
}
