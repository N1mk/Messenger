package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// CreateChat implements the create_chat endpoint on the business layer.
// The method receives dto.CreateChatInput and returns dto.CreateChatOutput and error.
// It creates new chat.
func (u *Usecase) CreateChat(input dto.CreateChatInput) (dto.CreateChatOutput, error) {
	added := make(map[int]struct{})

	creator, err := u.postgres.GetUserByID(input.CreatorID)
	if err != nil {
		return dto.CreateChatOutput{}, fmt.Errorf("user getting error: %w", err)
	}

	users := []domain.User{*creator}

	for _, username := range input.Usernames {
		user, err := u.postgres.GetUserByName(username)
		if err != nil {
			return dto.CreateChatOutput{}, fmt.Errorf("user getting error: %w", err)
		}

		_, exsists := added[user.ID]
		if exsists {
			continue
		}

		added[user.ID] = struct{}{}
		users = append(users, *user)
	}

	chat, err := u.postgres.CreateChat(domain.Chat{Name: input.Name, OwnerID: input.CreatorID, Users: users})
	if err != nil {
		return dto.CreateChatOutput{}, fmt.Errorf("chat creation error: %w", err)
	}

	creator.Chats = append(creator.Chats, *chat)

	err = u.postgres.UpdateUser(creator)
	if err != nil {
		return dto.CreateChatOutput{}, fmt.Errorf("user update error: %w", err)
	}

	return dto.CreateChatOutput{ChatID: chat.ID}, nil
}
