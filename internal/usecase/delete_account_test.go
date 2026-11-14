package usecase

import (
	"context"
	"errors"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"messenger/internal/usecase/mocks"
	"messenger/internal/usecase/test"
	"messenger/pkg/logger"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestDeleteAccount(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.DeleteAccountInput
		WantedOutput dto.DeleteAccountOutput
		WantedError  error

		RemoveUserWantedID       int
		RemoveUserReturningError error
	}{
		{
			Name:                     "ok",
			Input:                    dto.DeleteAccountInput{UserID: 1},
			WantedOutput:             dto.DeleteAccountOutput{},
			WantedError:              nil,
			RemoveUserWantedID:       1,
			RemoveUserReturningError: nil,
		},
		{
			Name:                     "error not found",
			Input:                    dto.DeleteAccountInput{UserID: 1},
			WantedOutput:             dto.DeleteAccountOutput{},
			WantedError:              domain.ErrNotFound,
			RemoveUserWantedID:       1,
			RemoveUserReturningError: domain.ErrNotFound,
		},
		{
			Name:                     "unknown error",
			Input:                    dto.DeleteAccountInput{UserID: 1},
			WantedOutput:             dto.DeleteAccountOutput{},
			WantedError:              test.UnknownTestErr,
			RemoveUserWantedID:       1,
			RemoveUserReturningError: test.UnknownTestErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(subT *testing.T) {
			ctrl := gomock.NewController(subT)
			rm := mocks.NewMockRepo(ctrl)
			tim := mocks.NewMockTokenIssuer(ctrl)

			rm.EXPECT().Begin(context.TODO())

			if tt.WantedError == nil {
				rm.EXPECT().Commit()
			} else {
				rm.EXPECT().Rollback()
			}

			rm.EXPECT().RemoveUser(tt.RemoveUserWantedID).Return(tt.RemoveUserReturningError)

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.DeleteAccount(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
