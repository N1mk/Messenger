package usecase

import (
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// UpdateAccount implements the update_account endpoint on the business layer.
// The method receives dto.UpdateAccountInput and returns dto.UpdateAccountOutput and error.
// It updates user with same id as in input.
func (u *Usecase) UpdateAccount(input dto.UpdateAccountInput) (dto.UpdateAccountOutput, error) {
	if err := u.postgres.UpdateUser(&domain.User{ID: input.UserID, Username: input.Username, Email: input.Email, Password: input.Password}); err != nil {
		return dto.UpdateAccountOutput{}, fmt.Errorf("user updating error: %w", err)
	}

	return dto.UpdateAccountOutput{}, nil
}
