package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"messenger/pkg/jwt"
)

// Login implements the login endpoint on the business layer.
// The method receives dto.LoginInput and returns dto.LoginOutput and error.
// It chceks if password is correct and returns user's ID and JWTToken.
// If password if incorrect, method returns domain.ErrIncorrectPassword.
func (u *Usecase) Login(input dto.LoginInput) (dto.LoginOutput, error) {
	user, err := u.postgres.GetUserByID(input.UserID)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("user get error: %w", err)
	}

	if user.Password != input.Password {
		return dto.LoginOutput{}, domain.ErrIncorrectPassword
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return dto.LoginOutput{}, fmt.Errorf("jswt-token generation error: %w", err)
	}

	return dto.LoginOutput{ID: user.ID, JWTToken: token}, nil
}
