package test

import (
	"errors"
	"messenger/internal/domain"
)

var UnknownTestErr = errors.New("unknown test error")

func CheckUserNameAndEmail(wantedUser domain.User) func(user domain.User) bool {
	return func(user domain.User) bool {
		if user.Username == wantedUser.Username && user.Email == wantedUser.Email {
			return true
		}

		return false
	}
}

func CheckUserIDAndNameAndEmail(wantedUser domain.User) func(user domain.User) bool {
	return func(user domain.User) bool {
		if user.ID == wantedUser.ID && user.Username == wantedUser.Username && user.Email == wantedUser.Email {
			return true
		}

		return false
	}
}

func CheckRefreshTokenUserID(wantedToken domain.RefreshToken) func(token domain.RefreshToken) bool {
	return func(token domain.RefreshToken) bool {
		if token.UserID == wantedToken.UserID {
			return true
		}

		return false
	}
}

func CheckMessageTextAndAndChatIDAndUserID(wantedMessage domain.Message) func(message domain.Message) bool {
	return func(message domain.Message) bool {
		if message.Text == wantedMessage.Text && message.ChatID == wantedMessage.ChatID && message.UserID == wantedMessage.UserID {
			return true
		}

		return false
	}
}
