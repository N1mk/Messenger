// Package usecase contains the buisnes logic of the application.
package usecase

import (
	"fmt"
	"messenger/internal/adapter/postgres"
	"messenger/internal/dto"
)

// UsecaseInterface defines all buisnes-logic methods.
type UsecaseInterface interface {
	Signup(input dto.SignupInput) (dto.SignupOutput, error)
	Login(input dto.LoginInput) (dto.LoginOutput, error)
	LoginByRefreshToken(input dto.LoginByRefreshTokenInput) (dto.LoginByRefreshTokenOutput, error)
	UpdateAccount(input dto.UpdateAccountInput) (dto.UpdateAccountOutput, error)
	DeleteAccount(input dto.DeleteAccountInput) (dto.DeleteAccountOutput, error)
	CreateChat(input dto.CreateChatInput) (dto.CreateChatOutput, error)
	UpdateChat(input dto.UpdateChatInput) (dto.UpdateChatOutput, error)
	DeleteChat(input dto.DeleteChatInput) (dto.DeleteChatOutput, error)
	GetChats(input dto.GetChatsInput) (dto.GetChatsOutput, error)
	SendMessage(input dto.SendMessageInput) (dto.SendMessageOutput, error)
	GetMessages(input dto.GetMessagesInput) (dto.GetMessagesOutput, error)
}

// Usecase is a default implementation of the UsecaseInterface.
type Usecase struct {
	postgres postgres.PostgresInterface
}

// NewUscase initializes Usecase and returns pointer on it.
func NewUsecase(postgres postgres.PostgresInterface) *Usecase {
	return &Usecase{postgres: postgres}
}

// isUserInChat is a privet Usecase method that checks if the user is in the chat.
// It returns (true, nil) if the user is in the chat, (false, nil) if the user is not in the chat and (false, err) if an error appeared while execution.
func (u *Usecase) isUserInChat(userID uint, chatID uint) (ok bool, err error) {
	chat, err := u.postgres.GetChatByID(chatID)
	if err != nil {
		return false, fmt.Errorf("chat getting error: %w", err)
	}

	for _, user := range chat.Users {
		if userID == user.ID {
			return true, nil
		}
	}

	return false, nil
}
