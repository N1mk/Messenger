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

func TestSignup(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.SignupInput
		WantedOutput dto.SignupOutput
		WantedError  error

		AddUserWantedUser     domain.User
		AddUserReturningUser  *domain.User
		AddUserReturningError error

		ExpectGenerateTokenPair                bool
		GenerateTokenPairWantedID              int
		GenerateTokenPairWantedName            string
		GenerateTokenPairReturningAccessToken  string
		GenerateTokenPairReturningRefreshToken string
		GenerateTokenPairReturningError        error

		ExpectAddRefreshToken         bool
		AddRefreshTokenWantedToken    domain.RefreshToken
		AddRefreshTokenReturningToken *domain.RefreshToken
		AddRefreshTokenReturningError error
	}{
		{
			Name:                                   "ok",
			Input:                                  dto.SignupInput{Username: "Nimk", Email: "nimk@mail.com", Password: "123456"},
			WantedOutput:                           dto.SignupOutput{ID: 1, AccesToken: "CoolAccessToken", RefreshToken: "CoolRefreshToken"},
			WantedError:                            nil,
			AddUserWantedUser:                      domain.User{Username: "Nimk", Email: "nimk@mail.com"},
			AddUserReturningUser:                   &domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")},
			AddUserReturningError:                  nil,
			ExpectGenerateTokenPair:                true,
			GenerateTokenPairWantedID:              1,
			GenerateTokenPairWantedName:            "Nimk",
			GenerateTokenPairReturningAccessToken:  "CoolAccessToken",
			GenerateTokenPairReturningRefreshToken: "CoolRefreshToken",
			GenerateTokenPairReturningError:        nil,
			ExpectAddRefreshToken:                  true,
			AddRefreshTokenWantedToken:             domain.RefreshToken{UserID: 1},
			AddRefreshTokenReturningToken:          &domain.RefreshToken{TokenHash: []byte("CoolTokenash"), UserID: 1},
			AddRefreshTokenReturningError:          nil,
		},
		{
			Name:                  "add user error duplicate",
			Input:                 dto.SignupInput{Username: "Nimk", Email: "nimk@mail.com", Password: "123456"},
			WantedOutput:          dto.SignupOutput{},
			WantedError:           domain.ErrDuplicated,
			AddUserWantedUser:     domain.User{Username: "Nimk", Email: "nimk@mail.com"},
			AddUserReturningUser:  nil,
			AddUserReturningError: domain.ErrDuplicated,
		},
		{
			Name:                  "add user unknown error",
			Input:                 dto.SignupInput{Username: "Nimk", Email: "nimk@mail.com", Password: "123456"},
			WantedOutput:          dto.SignupOutput{},
			WantedError:           test.UnknownTestErr,
			AddUserWantedUser:     domain.User{Username: "Nimk", Email: "nimk@mail.com"},
			AddUserReturningUser:  nil,
			AddUserReturningError: test.UnknownTestErr,
		},
		{
			Name:                                   "generate token pair unknown error",
			Input:                                  dto.SignupInput{Username: "Nimk", Email: "nimk@mail.com", Password: "123456"},
			WantedOutput:                           dto.SignupOutput{},
			WantedError:                            test.UnknownTestErr,
			AddUserWantedUser:                      domain.User{Username: "Nimk", Email: "nimk@mail.com"},
			AddUserReturningUser:                   &domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")},
			AddUserReturningError:                  nil,
			ExpectGenerateTokenPair:                true,
			GenerateTokenPairWantedID:              1,
			GenerateTokenPairWantedName:            "Nimk",
			GenerateTokenPairReturningAccessToken:  "CoolAccessToken",
			GenerateTokenPairReturningRefreshToken: "CoolRefreshToken",
			GenerateTokenPairReturningError:        test.UnknownTestErr,
		},
		{
			Name:                                   "add refresh token unknown error",
			Input:                                  dto.SignupInput{Username: "Nimk", Email: "nimk@mail.com", Password: "123456"},
			WantedOutput:                           dto.SignupOutput{},
			WantedError:                            test.UnknownTestErr,
			AddUserWantedUser:                      domain.User{Username: "Nimk", Email: "nimk@mail.com"},
			AddUserReturningUser:                   &domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")},
			AddUserReturningError:                  nil,
			ExpectGenerateTokenPair:                true,
			GenerateTokenPairWantedID:              1,
			GenerateTokenPairWantedName:            "Nimk",
			GenerateTokenPairReturningAccessToken:  "CoolAccessToken",
			GenerateTokenPairReturningRefreshToken: "CoolRefreshToken",
			GenerateTokenPairReturningError:        nil,
			ExpectAddRefreshToken:                  true,
			AddRefreshTokenWantedToken:             domain.RefreshToken{UserID: 1},
			AddRefreshTokenReturningToken:          &domain.RefreshToken{TokenHash: []byte("CoolTokenash"), UserID: 1},
			AddRefreshTokenReturningError:          test.UnknownTestErr,
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

			rm.EXPECT().AddUser(gomock.Cond(test.CheckUserNameAndEmail(tt.AddUserWantedUser))).Return(tt.AddUserReturningUser, tt.AddUserReturningError)

			if tt.ExpectGenerateTokenPair {
				tim.EXPECT().GenerateTokenPair(tt.GenerateTokenPairWantedID, tt.GenerateTokenPairWantedName).Return(tt.GenerateTokenPairReturningAccessToken, tt.GenerateTokenPairReturningRefreshToken, tt.GenerateTokenPairReturningError)
			}

			if tt.ExpectAddRefreshToken {
				rm.EXPECT().AddRefreshToken(gomock.Cond(test.CheckRefreshTokenUserID(tt.AddRefreshTokenWantedToken))).Return(tt.AddRefreshTokenReturningToken, tt.AddRefreshTokenReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.Signup(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
