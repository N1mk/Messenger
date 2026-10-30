package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"

	"golang.org/x/crypto/bcrypt"
)

// UpdateAccount implements the update_account endpoint on the business layer.
// The method receives dto.UpdateAccountInput and returns dto.UpdateAccountOutput and error.
// It updates user with same id as in input.
func (u *Usecase) UpdateAccount(input dto.UpdateAccountInput) (dto.UpdateAccountOutput, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UpdateAccountOutput{}, fmt.Errorf("password hash generation error: %w", err)
	}

	if err := u.postgres.UpdateUser(&domain.User{ID: input.UserID, Username: input.Username, Email: input.Email, PasswordHash: passwordHash}); err != nil {
		return dto.UpdateAccountOutput{}, fmt.Errorf("user updating error: %w", err)
	}

	return dto.UpdateAccountOutput{}, nil
}
