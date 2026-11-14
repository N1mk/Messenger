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

func TestLogin(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.LoginInput
		WantedOutput dto.LoginOutput
		WantedError  error

		GetUserByIDWantedID       int
		GetUserByIDReturningUser  *domain.User
		GetUserByIDReturningError error

		ExpectGenerateTokenPair                bool
		GenerateTokenPairWantedID              int
		GenerateTokenPairWantedUsername        string
		GenerateTokenPairReturningAccessToken  string
		GenerateTokenPairReturningRefreshToken string
		GenerateTokenPairReturningError        error

		ExpectReplaceRefreshToken         bool
		ReplaceRefreshTokenWantedToken    domain.RefreshToken
		ReplaceRefreshTokenReturningError error
	}{
		{
			Name:                                   "ok",
			Input:                                  dto.LoginInput{UserID: 1, Password: "123456"},
			WantedOutput:                           dto.LoginOutput{AccesToken: "CoolNewAccessToken", RefreshToken: "CoolNewRefreshToken"},
			WantedError:                            nil,
			GetUserByIDWantedID:                    1,
			GetUserByIDReturningUser:               &domain.User{ID: 1, Username: "Nimk", PasswordHash: []byte("$2a$10$Ng8RFCyV/E6WeTyTz6FdhOc.Y6jyTENt67T6D6ViLtGq27k57rzC6")},
			GetUserByIDReturningError:              nil,
			ExpectGenerateTokenPair:                true,
			GenerateTokenPairWantedID:              1,
			GenerateTokenPairWantedUsername:        "Nimk",
			GenerateTokenPairReturningAccessToken:  "CoolNewAccessToken",
			GenerateTokenPairReturningRefreshToken: "CoolNewRefreshToken",
			GenerateTokenPairReturningError:        nil,
			ExpectReplaceRefreshToken:              true,
			ReplaceRefreshTokenWantedToken:         domain.RefreshToken{UserID: 1},
			ReplaceRefreshTokenReturningError:      nil,
		},
		{
			Name:                      "get user by id error not found",
			Input:                     dto.LoginInput{UserID: 1, Password: "123456"},
			WantedOutput:              dto.LoginOutput{},
			WantedError:               domain.ErrNotFound,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  nil,
			GetUserByIDReturningError: domain.ErrNotFound,
		},
		{
			Name:                      "get user by id unknown error",
			Input:                     dto.LoginInput{UserID: 1, Password: "123456"},
			WantedOutput:              dto.LoginOutput{},
			WantedError:               test.UnknownTestErr,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  nil,
			GetUserByIDReturningError: test.UnknownTestErr,
		},
		{
			Name:                                   "generate token pair unknown error",
			Input:                                  dto.LoginInput{UserID: 1, Password: "123456"},
			WantedOutput:                           dto.LoginOutput{},
			WantedError:                            test.UnknownTestErr,
			GetUserByIDWantedID:                    1,
			GetUserByIDReturningUser:               &domain.User{ID: 1, Username: "Nimk", PasswordHash: []byte("$2a$10$Ng8RFCyV/E6WeTyTz6FdhOc.Y6jyTENt67T6D6ViLtGq27k57rzC6")},
			GetUserByIDReturningError:              nil,
			ExpectGenerateTokenPair:                true,
			GenerateTokenPairWantedID:              1,
			GenerateTokenPairWantedUsername:        "Nimk",
			GenerateTokenPairReturningAccessToken:  "",
			GenerateTokenPairReturningRefreshToken: "",
			GenerateTokenPairReturningError:        test.UnknownTestErr,
		},
		{
			Name:                                   "replace refresh token error not found",
			Input:                                  dto.LoginInput{UserID: 1, Password: "123456"},
			WantedOutput:                           dto.LoginOutput{},
			WantedError:                            domain.ErrNotFound,
			GetUserByIDWantedID:                    1,
			GetUserByIDReturningUser:               &domain.User{ID: 1, Username: "Nimk", PasswordHash: []byte("$2a$10$Ng8RFCyV/E6WeTyTz6FdhOc.Y6jyTENt67T6D6ViLtGq27k57rzC6")},
			GetUserByIDReturningError:              nil,
			ExpectGenerateTokenPair:                true,
			GenerateTokenPairWantedID:              1,
			GenerateTokenPairWantedUsername:        "Nimk",
			GenerateTokenPairReturningAccessToken:  "CoolNewAccessToken",
			GenerateTokenPairReturningRefreshToken: "CoolNewRefreshToken",
			GenerateTokenPairReturningError:        nil,
			ExpectReplaceRefreshToken:              true,
			ReplaceRefreshTokenWantedToken:         domain.RefreshToken{UserID: 1},
			ReplaceRefreshTokenReturningError:      domain.ErrNotFound,
		},
		{
			Name:                                   "replace refresh token unknown error",
			Input:                                  dto.LoginInput{UserID: 1, Password: "123456"},
			WantedOutput:                           dto.LoginOutput{},
			WantedError:                            test.UnknownTestErr,
			GetUserByIDWantedID:                    1,
			GetUserByIDReturningUser:               &domain.User{ID: 1, Username: "Nimk", PasswordHash: []byte("$2a$10$Ng8RFCyV/E6WeTyTz6FdhOc.Y6jyTENt67T6D6ViLtGq27k57rzC6")},
			GetUserByIDReturningError:              nil,
			ExpectGenerateTokenPair:                true,
			GenerateTokenPairWantedID:              1,
			GenerateTokenPairWantedUsername:        "Nimk",
			GenerateTokenPairReturningAccessToken:  "CoolNewAccessToken",
			GenerateTokenPairReturningRefreshToken: "CoolNewRefreshToken",
			GenerateTokenPairReturningError:        nil,
			ExpectReplaceRefreshToken:              true,
			ReplaceRefreshTokenWantedToken:         domain.RefreshToken{UserID: 1},
			ReplaceRefreshTokenReturningError:      test.UnknownTestErr,
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

			rm.EXPECT().GetUserByID(tt.GetUserByIDWantedID).Return(tt.GetUserByIDReturningUser, tt.GetUserByIDReturningError)

			if tt.ExpectGenerateTokenPair {
				tim.EXPECT().GenerateTokenPair(tt.GenerateTokenPairWantedID, tt.GenerateTokenPairWantedUsername).Return(tt.GenerateTokenPairReturningAccessToken, tt.GenerateTokenPairReturningRefreshToken, tt.GenerateTokenPairReturningError)
			}

			if tt.ExpectReplaceRefreshToken {
				rm.EXPECT().ReplaceRefreshToken(gomock.Cond(test.CheckRefreshTokenUserID(tt.ReplaceRefreshTokenWantedToken))).Return(tt.ReplaceRefreshTokenReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.Login(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
