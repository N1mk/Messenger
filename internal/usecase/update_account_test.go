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

func TestUpdateAccount(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.UpdateAccountInput
		WantedOutput dto.UpdateAccountOutput
		WantedError  error

		UpdateUserWantedUser     domain.User
		UpdateUserReturningError error
	}{
		{
			Name:                     "ok",
			Input:                    dto.UpdateAccountInput{UserID: 1, Username: "Nimk", Email: "nimk@mail.com", Password: "123456"},
			WantedOutput:             dto.UpdateAccountOutput{},
			WantedError:              nil,
			UpdateUserWantedUser:     domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com"},
			UpdateUserReturningError: nil,
		},
		{
			Name:                     "update user unknown error",
			Input:                    dto.UpdateAccountInput{UserID: 1, Username: "Nimk", Email: "nimk@mail.com", Password: "123456"},
			WantedOutput:             dto.UpdateAccountOutput{},
			WantedError:              test.UnknownTestErr,
			UpdateUserWantedUser:     domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com"},
			UpdateUserReturningError: test.UnknownTestErr,
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

			rm.EXPECT().UpdateUser(gomock.Cond(test.CheckUserIDAndNameAndEmail(tt.UpdateUserWantedUser))).Return(tt.UpdateUserReturningError)

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.UpdateAccount(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
