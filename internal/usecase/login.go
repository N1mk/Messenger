package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"messenger/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

// Login implements the login endpoint on the business layer.
// The method receives dto.LoginInput and returns dto.LoginOutput and error.
// It chceks if password is correct and returns user's ID and JWT-Token pair.
// If password if incorrect, method returns domain.ErrIncorrectPassword.
func (u *Usecase) Login(input dto.LoginInput) (dto.LoginOutput, error) {
	user, err := u.postgres.GetUserByID(input.UserID)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("user get error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(input.Password)); err != nil {
		return dto.LoginOutput{}, fmt.Errorf("%w: password hash comparing error: %w", domain.ErrUnauthorized, err)
	}

	tokenPair, err := jwt.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("jwt-token pair generation error: %w", err)
	}

	if err := u.postgres.ReplaceRefreshToken(tokenPair.RefreshToken, input.UserID); err != nil {
		return dto.LoginOutput{}, fmt.Errorf("refresh token replacing error: %w", err)
	}

	return dto.LoginOutput{ID: user.ID, TokenPair: tokenPair}, nil
}
