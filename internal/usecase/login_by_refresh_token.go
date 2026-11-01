package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// LoginByRefreshToken implements the login_by_refresh_token endpoint on the business layer.
// The method receives dto.LoginByRefreshTokenInput and returns dto.LoginByRefreshTokenOutput and error.
// It chceks if refresh token exists and valid and returns user's ID and JWT-Token pair.
// If refresh token is not exist, method returns domain.ErrUnauthorized.
func (u *Usecase) LoginByRefreshToken(ctx context.Context, input dto.LoginByRefreshTokenInput) (dto.LoginByRefreshTokenOutput, error) {
	user, err := u.p.GetUserByID(ctx, input.UserID)
	if err != nil {
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("user get error: %w", err)
	}

	refreshTokenHash := sha256.Sum256([]byte(input.RefreshToken))

	if ok, err := u.p.RefreshTokenExists(ctx, refreshTokenHash[:]); !ok {
		if err != nil {
			return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("refresh token exist check error: %w", err)
		}
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("refresh token doesn't exist: %w", domain.ErrUnauthorized)
	}

	userID, err := u.ti.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("refresh token parsing error: %w", err)
	}

	if input.UserID != userID {
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("refresh token user id and request user id don't match: %s", domain.ErrUnauthorized)
	}

	accesToken, refreshToken, err := u.ti.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("jwt-token pair generation error: %w", err)
	}

	newRefreshTokenHash := sha256.Sum256([]byte(refreshToken))

	if err := u.p.ReplaceRefreshToken(ctx, newRefreshTokenHash[:], input.UserID); err != nil {
		return dto.LoginByRefreshTokenOutput{}, fmt.Errorf("refresh token replacing error: %w", err)
	}

	return dto.LoginByRefreshTokenOutput{ID: user.ID, AccesToken: accesToken, RefreshToken: refreshToken}, nil
}
