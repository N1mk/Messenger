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
	u.r.Begin(ctx)
	defer func() {
		if err != nil {
			u.r.Rollback()
			return
		}
		u.r.Commit()
	}()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		u.l.LogError(fmt.Sprintf("password hash generation error: %s", err.Error()))
		return dto.SignupOutput{}, fmt.Errorf("password hash generation error: %w", err)
	}

	user := domain.User{Username: input.Username, PasswordHash: passwordHash, Email: input.Email}

	newUser, err := u.r.CreateUser(user)
	if err != nil {
		u.l.LogError(fmt.Sprintf("user creation error: %s", err.Error()))
		return dto.SignupOutput{}, fmt.Errorf("user creation error: %w", err)
	}

	accesToken, refreshToken, err := u.ti.GenerateTokenPair(newUser.ID, newUser.Username)
	if err != nil {
		u.l.LogError(fmt.Sprintf("jwt-token generation error: %s", err.Error()))
		return dto.SignupOutput{}, fmt.Errorf("jwt-token generation error: %w", err)
	}

	refreshTokenHash := sha256.Sum256([]byte(refreshToken))

	if err := u.r.CreateRefreshToken(refreshTokenHash[:], newUser.ID); err != nil {
		u.l.LogError(fmt.Sprintf("refresh token creating error: %s", err.Error()))
		return dto.SignupOutput{}, fmt.Errorf("refresh token creating error: %w", err)
	}

	return dto.SignupOutput{ID: newUser.ID, AccesToken: accesToken, RefreshToken: refreshToken}, nil
}
