package usecase

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"

	"golang.org/x/crypto/bcrypt"
)

// UpdateAccount implements the update_account endpoint on the business layer.
// The method receives dto.UpdateAccountInput and returns dto.UpdateAccountOutput and error.
// It updates user with same id as in input.
func (u *Usecase) UpdateAccount(ctx context.Context, input dto.UpdateAccountInput) (output dto.UpdateAccountOutput, err error) {
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
		return dto.UpdateAccountOutput{}, fmt.Errorf("password hash generation error: %w", err)
	}

	if err := u.r.UpdateUser(&domain.User{ID: input.UserID, Username: input.Username, Email: input.Email, PasswordHash: passwordHash}); err != nil {
		u.l.LogError(fmt.Sprintf("user updating error: %s", err.Error()))
		return dto.UpdateAccountOutput{}, fmt.Errorf("user updating error: %w", err)
	}

	return dto.UpdateAccountOutput{}, nil
}
