// Package usecase contains the buisnes logic of the application.
package usecase

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// UsecaseInterface defines all buisnes-logic methods.
type UsecaseInterface interface {
	// Signup implements the signup endpoint on the business layer.
	// The method receives dto.SignupInput and context.Context and returns dto.SignupOutput and error.
	// It creates new user from input and rutruns user's ID and JWT-Token pair.
	Signup(ctx context.Context, input dto.SignupInput) (dto.SignupOutput, error)

	// Login implements the login endpoint on the business layer.
	// The method receives dto.LoginInput and context.Context and returns dto.LoginOutput and error.
	// It chceks if password is correct and returns user's ID and JWT-Token pair.
	// If password if incorrect, method returns domain.ErrUnauthorized.
	Login(ctx context.Context, input dto.LoginInput) (dto.LoginOutput, error)

	// LoginByRefreshToken implements the login_by_refresh_token endpoint on the business layer.
	// The method receives dto.LoginByRefreshTokenInput and context.Context and returns dto.LoginByRefreshTokenOutput and error.
	// It chceks if refresh token exists and valid and returns user's ID and JWT-Token pair.
	// If refresh token is not exist, method returns domain.ErrUnauthorized.
	LoginByRefreshToken(ctx context.Context, input dto.LoginByRefreshTokenInput) (dto.LoginByRefreshTokenOutput, error)

	// UpdateAccount implements the update_account endpoint on the business layer.
	// The method receives dto.UpdateAccountInput and context.Context and returns dto.UpdateAccountOutput and error.
	// It updates user with same id as in input.
	UpdateAccount(ctx context.Context, input dto.UpdateAccountInput) (dto.UpdateAccountOutput, error)

	// DeleteAccount implements the delete_account endpoint on the business layer.
	// The method receives dto.DeleteAccountInput and context.Context and returns dto.DeleteAccountOutput and error.
	// It deletes user by id.
	DeleteAccount(ctx context.Context, input dto.DeleteAccountInput) (dto.DeleteAccountOutput, error)

	// CreateChat implements the create_chat endpoint on the business layer.
	// The method receives dto.CreateChatInput and context.Context and returns dto.CreateChatOutput and error.
	// It creates new chat.
	CreateChat(ctx context.Context, input dto.CreateChatInput) (dto.CreateChatOutput, error)

	// UpdateChat implements the update_chat endpoint on the business layer.
	// The method receives dto.UpdateChatInput and context.Context and returns dto.UpdateChatOutput and error.
	// It updates chat with same id as in input.
	UpdateChat(ctx context.Context, input dto.UpdateChatInput) (dto.UpdateChatOutput, error)

	// DeleteChat implements the delete_chat endpoint on the business layer.
	// The method receives dto.DeleteChatInput and context.Context and returns dto.DeleteChatOutput and error.
	// It deletes chat by id.
	// If uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
	DeleteChat(ctx context.Context, input dto.DeleteChatInput) (dto.DeleteChatOutput, error)

	// GetChats implements the get_chats endpoint on the business layer.
	// The method receives dto.GetChatsInput and context.Context and returns dto.GetChatsOutput and error.
	// It returns all chats, in which the user is membered.
	GetChats(ctx context.Context, input dto.GetChatsInput) (dto.GetChatsOutput, error)

	// SendMessage implements the send_message endpoint on the business layer.
	// The method receives dto.SendMessageInput and context.Context and returns dto.SendMessageOutput and error.
	// It creates new message in chat with ID = input.ChatID.
	// If the uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
	SendMessage(ctx context.Context, input dto.SendMessageInput) (dto.SendMessageOutput, error)

	// GetMessages implements the get_messages endpoint on the business layer.
	// The method receives dto.GetMessagesInput and context.Context and returns dto.GetMessagesOutput and error.
	// It returns all messages from chat.
	// If the uesr is not membered in the chat, method returns dto.ErrUserIsNotInChat.
	GetMessages(ctx context.Context, input dto.GetMessagesInput) (dto.GetMessagesOutput, error)
}

// Repo defines all methods that usecase need from repo.
type Repo interface {
	// Begins transaction
	Begin(ctx context.Context)
	// Rollbacks transaction
	Rollback()
	// Commits transaction
	Commit()
	// GetUserByID selects user from the users table by id.
	GetUserByID(id int) (*domain.User, error)

	// GetUserByName selects user from the users table by name.
	GetUserByName(username string) (*domain.User, error)

	// СreateUser creates new element in the users table.
	CreateUser(user domain.User) (*domain.User, error)

	// RemoveUser removes user from the users table by id.
	RemoveUser(id int) error

	// UpdateUser updates user with correct id.
	UpdateUser(user *domain.User) error

	// GetChatByID selects chat from the chats table by id.
	GetChatByID(id int) (*domain.Chat, error)

	// СreateChat creates new element in the chats table.
	CreateChat(chat domain.Chat) (*domain.Chat, error)

	// RemoveChat removes chat from the chats table by id.
	RemoveChat(id int) error

	// UpdateChat updates chat with correct id.
	UpdateChat(chat *domain.Chat) error

	// GetMessages selects all messages from the messages table with correct chat id.
	GetMessages(chatID int, lastID int, pageSize int) ([]*domain.Message, error)

	// RefreshTokenExists checks if refresh token exists in the refresh_tokens table.
	RefreshTokenExists(tokenHash []byte) (ok bool, err error)

	// СreateRefreshToken creates new element in the refresh_tokens table.
	CreateRefreshToken(tokenHash []byte, userID int) error

	// RemoveRefreshToken removes refresh token from refresh_tokens table.
	RemoveRefreshToken(userID int) error

	// ReplaceRefreshToken replaces the old user refresh token with the new one in transaction in the refresh_tokens table.
	ReplaceRefreshToken(newtokenHash []byte, userID int) error
}

// TokenIssuer is a interface which defines work with JWT tokens
type TokenIssuer interface {
	GenerateTokenPair(id int, username string) (accesToken string, refreshToken string, err error)
	ParseAccesToken(accesToken string) (id int, username string, err error)
	ParseRefreshToken(refreshToken string) (id int, err error)
}

// Usecase is a default implementation of the UsecaseInterface.
type Usecase struct {
	r  Repo
	ti TokenIssuer
}

// NewUscase initializes Usecase and returns pointer on it.
func NewUsecase(r Repo, ti TokenIssuer) *Usecase {
	return &Usecase{r: r, ti: ti}
}

// isUserInChat is a privet Usecase method that checks if the user is in the chat.
// It returns (true, nil) if the user is in the chat, (false, nil) if the user is not in the chat and (false, err) if an error appeared while execution.
func (u *Usecase) isUserInChat(userID int, chatID int) (ok bool, err error) {
	chat, err := u.r.GetChatByID(chatID)
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
