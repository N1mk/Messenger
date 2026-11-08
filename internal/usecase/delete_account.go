package usecase

import (
	"context"
	"fmt"
	"messenger/internal/dto"
)

// DeleteAccount implements the delete_account endpoint on the business layer.
// The method receives dto.DeleteAccountInput and returns dto.DeleteAccountOutput and error.
// It deletes user by id.
func (u *Usecase) DeleteAccount(ctx context.Context, input dto.DeleteAccountInput) (output dto.DeleteAccountOutput, err error) {
	u.r.Begin(ctx)
	defer func() {
		if err != nil {
			u.r.Rollback()
			return
		}
		u.r.Commit()
	}()

	if err := u.r.RemoveUser(input.UserID); err != nil {
		u.l.LogError(fmt.Sprintf("user removing error: %s", err.Error()))
		return dto.DeleteAccountOutput{}, fmt.Errorf("user removing error: %w", err)
	}

	return dto.DeleteAccountOutput{}, nil
}
