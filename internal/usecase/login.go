package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"

	"golang.org/x/crypto/bcrypt"
)

// Login implements the login endpoint on the business layer.
// The method receives dto.LoginInput and returns dto.LoginOutput and error.
// It chceks if password is correct and returns user's ID and JWT-Token pair.
// If password if incorrect, method returns domain.ErrUnauthorized.
func (u *Usecase) Login(ctx context.Context, input dto.LoginInput) (output dto.LoginOutput, err error) {
	u.p.Begin(ctx)
	defer func() {
		if err != nil {
			u.p.Rollback()
			return
		}
		u.p.Commit()
	}()

	user, err := u.p.GetUserByID(input.UserID)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("user get error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(input.Password)); err != nil {
		return dto.LoginOutput{}, fmt.Errorf("%w: password hash comparing error: %w", domain.ErrUnauthorized, err)
	}

	accesToken, refreshToken, err := u.ti.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("jwt-token pair generation error: %w", err)
	}

	refreshTokenHash := sha256.Sum256([]byte(refreshToken))
	if err := u.p.ReplaceRefreshToken(refreshTokenHash[:], input.UserID); err != nil {
		return dto.LoginOutput{}, fmt.Errorf("refresh token replacing error: %w", err)
	}

	return dto.LoginOutput{ID: user.ID, AccesToken: accesToken, RefreshToken: refreshToken}, nil
}
