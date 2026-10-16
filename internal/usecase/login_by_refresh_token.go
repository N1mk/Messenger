package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"messenger/pkg/jwt"
)

// LoginByRefreshToken implements the login_by_refresh_token endpoint on the business layer.
// The method receives dto.LoginByRefreshTokenInput and returns dto.LoginByRefreshTokenOutput and error.
// It chceks if refresh token exists and valid and returns user's ID and JWT-Token pair.
// If refresh token is not exist, method returns domain.ErrUnauthorized.
func (u *Usecase) LoginByRefreshToken(input dto.LoginByRefreshTokenInput) (dto.LoginByRefreshTokenOutput, error) {
	user, err := u.postgres.GetUserByID(input.UserID)
	if err != nil {
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("user get error: %w", err)
	}

	if ok, err := u.postgres.RefreshTokenExists(jwt.RefreshToken{Token: input.RefreshToken}); !ok {
		if err != nil {
			return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("refresh token exist check error: %w", err)
		}
		return dto.LoginByRefreshTokenOutput{}, domain.ErrUnauthorized
	}

	if err := jwt.ParseRefreshToken(input.RefreshToken); err != nil {
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("refresh token parsing error: %w", err)
	}

	tokenPair, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("jwt-token pair generation error: %w", err)
	}

	if err := u.postgres.CreateRefreshToken(jwt.RefreshToken{Token: tokenPair.RefreshToken}); err != nil {
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("refresh token adding to database error: %w", err)
	}

	return dto.LoginByRefreshTokenOutput{ID: user.ID, TokenPair: tokenPair}, nil
}
