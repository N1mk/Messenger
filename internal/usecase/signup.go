package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"

	"golang.org/x/crypto/bcrypt"
)

// Signup implements the signup endpoint on the business layer.
// The method receives dto.SignupInput and returns dto.SignupOutput and error.
// It creates new user from input and rutruns user's ID and JWT-Token pair.
func (u *Usecase) Signup(ctx context.Context, input dto.SignupInput) (output dto.SignupOutput, err error) {
	u.p.Begin(ctx)
	defer func() {
		if err != nil {
			u.p.Rollback()
			return
		}
		u.p.Commit()
	}()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.SignupOutput{}, fmt.Errorf("password hash generation error: %w", err)
	}

	user := domain.User{Username: input.Username, PasswordHash: passwordHash, Email: input.Email}

	newUser, err := u.p.CreateUser(user)
	if err != nil {
		return dto.SignupOutput{}, fmt.Errorf("user creation error: %w", err)
	}

	accesToken, refreshToken, err := u.ti.GenerateTokenPair(newUser.ID, newUser.Username)
	if err != nil {
		return dto.SignupOutput{}, fmt.Errorf("jswt-token generation error: %w", err)
	}

	refreshTokenHash := sha256.Sum256([]byte(refreshToken))

	if err := u.p.ReplaceRefreshToken(refreshTokenHash[:], newUser.ID); err != nil {
		return dto.SignupOutput{}, fmt.Errorf("refresh token replacing error: %w", err)
	}

	return dto.SignupOutput{ID: newUser.ID, AccesToken: accesToken, RefreshToken: refreshToken}, nil
}
