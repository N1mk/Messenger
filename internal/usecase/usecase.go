// Package usecase contains the buisnes logic of the application.
package usecase

import (
	"fmt"
	"messenger/internal/adapter/postgres"
	"messenger/internal/dto"
)

// UsecaseInterface defines all buisnes-logic methods.
type UsecaseInterface interface {
	// Signup implements the signup endpoint on the business layer.
	// The method receives dto.SignupInput and returns dto.SignupOutput and error.
	// It creates new user from input and rutruns user's ID and JWT-Token pair.
	Signup(input dto.SignupInput) (dto.SignupOutput, error)

	// Login implements the login endpoint on the business layer.
	// The method receives dto.LoginInput and returns dto.LoginOutput and error.
	// It chceks if password is correct and returns user's ID and JWT-Token pair.
	// If password if incorrect, method returns domain.ErrIncorrectPassword.
	Login(input dto.LoginInput) (dto.LoginOutput, error)

	// LoginByRefreshToken implements the login_by_refresh_token endpoint on the business layer.
	// The method receives dto.LoginByRefreshTokenInput and returns dto.LoginByRefreshTokenOutput and error.
	// It chceks if refresh token exists and valid and returns user's ID and JWT-Token pair.
	// If refresh token is not exist, method returns domain.ErrUnauthorized.
	LoginByRefreshToken(input dto.LoginByRefreshTokenInput) (dto.LoginByRefreshTokenOutput, error)

	// UpdateAccount implements the update_account endpoint on the business layer.
	// The method receives dto.UpdateAccountInput and returns dto.UpdateAccountOutput and error.
	// It updates user with same id as in input.
	UpdateAccount(input dto.UpdateAccountInput) (dto.UpdateAccountOutput, error)

	// DeleteAccount implements the delete_account endpoint on the business layer.
	// The method receives dto.DeleteAccountInput and returns dto.DeleteAccountOutput and error.
	// It deletes user by id.
	DeleteAccount(input dto.DeleteAccountInput) (dto.DeleteAccountOutput, error)

	// CreateChat implements the create_chat endpoint on the business layer.
	// The method receives dto.CreateChatInput and returns dto.CreateChatOutput and error.
	// It creates new chat.
	CreateChat(input dto.CreateChatInput) (dto.CreateChatOutput, error)

	// UpdateChat implements the update_chat endpoint on the business layer.
	// The method receives dto.UpdateChatInput and returns dto.UpdateChatOutput and error.
	// It updates chat with same id as in input.
	UpdateChat(input dto.UpdateChatInput) (dto.UpdateChatOutput, error)

	// DeleteChat implements the delete_chat endpoint on the business layer.
	// The method receives dto.DeleteChatInput and returns dto.DeleteChatOutput and error.
	// It deletes chat by id.
	// If uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
	DeleteChat(input dto.DeleteChatInput) (dto.DeleteChatOutput, error)

	// GetChats implements the get_chats endpoint on the business layer.
	// The method receives dto.GetChatsInput and returns dto.GetChatsOutput and error.
	// It returns all chats, in which the user is membered.
	GetChats(input dto.GetChatsInput) (dto.GetChatsOutput, error)

	// SendMessage implements the send_message endpoint on the business layer.
	// The method receives dto.SendMessageInput and returns dto.SendMessageOutput and error.
	// It creates new message in chat with ID = input.ChatID.
	// If the uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
	SendMessage(input dto.SendMessageInput) (dto.SendMessageOutput, error)

	// GetMessages implements the get_messages endpoint on the business layer.
	// The method receives dto.GetMessagesInput and returns dto.GetMessagesOutput and error.
	// It returns all messages from chat.
	// If the uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
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
func (u *Usecase) isUserInChat(userID int, chatID int) (ok bool, err error) {
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
