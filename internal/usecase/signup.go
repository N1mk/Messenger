package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"messenger/pkg/jwt"
)

// Signup implements the signup endpoint on the business layer.
// The method receives dto.SignupInput and returns dto.SignupOutput and error.
// It creates new user from input and rutruns user's ID and JWT-Token pair.
func (u *Usecase) Signup(input dto.SignupInput) (dto.SignupOutput, error) {
	user := domain.User{Username: input.Username, Password: input.Password, Email: input.Email}

	newUser, err := u.postgres.CreateUser(user)
	if err != nil {
		return dto.SignupOutput{}, fmt.Errorf("user creation error: %w", err)
	}

	tokenPair, err := jwt.GenerateToken(newUser.ID, newUser.Username)
	if err != nil {
		return dto.SignupOutput{}, fmt.Errorf("jswt-token generation error: %w", err)
	}

	if err := u.postgres.CreateRefreshToken(jwt.RefreshToken{Token: tokenPair.RefreshToken}); err != nil {
		return dto.SignupOutput{}, fmt.Errorf("refresh token adding to database error: %w", err)
	}

	return dto.SignupOutput{ID: newUser.ID, TokenPair: tokenPair}, nil
}
