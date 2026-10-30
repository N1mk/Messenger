package usecase

import (
	"crypto/sha256"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"messenger/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

// Signup implements the signup endpoint on the business layer.
// The method receives dto.SignupInput and returns dto.SignupOutput and error.
// It creates new user from input and rutruns user's ID and JWT-Token pair.
func (u *Usecase) Signup(input dto.SignupInput) (dto.SignupOutput, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.SignupOutput{}, fmt.Errorf("password hash generation error: %w", err)
	}

	user := domain.User{Username: input.Username, PasswordHash: passwordHash, Email: input.Email}

	newUser, err := u.postgres.CreateUser(user)
	if err != nil {
		return dto.SignupOutput{}, fmt.Errorf("user creation error: %w", err)
	}

	tokenPair, err := jwt.GenerateTokenPair(newUser.ID, newUser.Username)
	if err != nil {
		return dto.SignupOutput{}, fmt.Errorf("jswt-token generation error: %w", err)
	}

	refreshTokenHash := sha256.Sum256([]byte(tokenPair.RefreshToken))

	if err := u.postgres.ReplaceRefreshToken(refreshTokenHash[:], newUser.ID); err != nil {
		return dto.SignupOutput{}, fmt.Errorf("refresh token replacing error: %w", err)
	}

	return dto.SignupOutput{ID: newUser.ID, TokenPair: tokenPair}, nil
}
