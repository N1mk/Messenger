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

func TestLoginByRefreshToken(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.LoginByRefreshTokenInput
		WantedOutput dto.LoginByRefreshTokenOutput
		WantedError  error

		GetUserByIDWantedID       int
		GetUserByIDReturningUser  *domain.User
		GetUserByIDReturningError error

		ExpectRefreshTokenExists         bool
		RefreshTokenExistsWantedToken    string
		RefreshTokenExistsReturningBool  bool
		RefreshTokenExistsReturningError error

		ExpectParseRefreshToken         bool
		ParseRefreshTokenWantedToken    string
		ParseRefreshTokenReturningID    int
		ParseRefreshTokenReturningError error

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
			Input:                                  dto.LoginByRefreshTokenInput{UserID: 1, RefreshToken: "CoolRefreshToken"},
			WantedOutput:                           dto.LoginByRefreshTokenOutput{AccesToken: "CoolNewAccessToken", RefreshToken: "CoolNewRefreshToken"},
			WantedError:                            nil,
			GetUserByIDWantedID:                    1,
			GetUserByIDReturningUser:               &domain.User{ID: 1, Username: "Nimk"},
			GetUserByIDReturningError:              nil,
			ExpectRefreshTokenExists:               true,
			RefreshTokenExistsWantedToken:          "CoolRefreshToken",
			RefreshTokenExistsReturningBool:        true,
			RefreshTokenExistsReturningError:       nil,
			ExpectParseRefreshToken:                true,
			ParseRefreshTokenWantedToken:           "CoolRefreshToken",
			ParseRefreshTokenReturningID:           1,
			ParseRefreshTokenReturningError:        nil,
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
			Input:                     dto.LoginByRefreshTokenInput{UserID: 1, RefreshToken: "CoolRefreshToken"},
			WantedOutput:              dto.LoginByRefreshTokenOutput{},
			WantedError:               domain.ErrNotFound,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  nil,
			GetUserByIDReturningError: domain.ErrNotFound,
		},
		{
			Name:                      "get user by id unknown error",
			Input:                     dto.LoginByRefreshTokenInput{UserID: 1, RefreshToken: "CoolRefreshToken"},
			WantedOutput:              dto.LoginByRefreshTokenOutput{},
			WantedError:               test.UnknownTestErr,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  nil,
			GetUserByIDReturningError: test.UnknownTestErr,
		},
		{
			Name:                             "refresh token exists unknown error",
			Input:                            dto.LoginByRefreshTokenInput{UserID: 1, RefreshToken: "CoolRefreshToken"},
			WantedOutput:                     dto.LoginByRefreshTokenOutput{},
			WantedError:                      test.UnknownTestErr,
			GetUserByIDWantedID:              1,
			GetUserByIDReturningUser:         &domain.User{ID: 1, Username: "Nimk"},
			GetUserByIDReturningError:        nil,
			ExpectRefreshTokenExists:         true,
			RefreshTokenExistsWantedToken:    "CoolRefreshToken",
			RefreshTokenExistsReturningBool:  false,
			RefreshTokenExistsReturningError: test.UnknownTestErr,
		},
		{
			Name:                             "refresh token not found",
			Input:                            dto.LoginByRefreshTokenInput{UserID: 1, RefreshToken: "CoolRefreshToken"},
			WantedOutput:                     dto.LoginByRefreshTokenOutput{},
			WantedError:                      domain.ErrUnauthorized,
			GetUserByIDWantedID:              1,
			GetUserByIDReturningUser:         &domain.User{ID: 1, Username: "Nimk"},
			GetUserByIDReturningError:        nil,
			ExpectRefreshTokenExists:         true,
			RefreshTokenExistsWantedToken:    "CoolRefreshToken",
			RefreshTokenExistsReturningBool:  false,
			RefreshTokenExistsReturningError: nil,
		},
		{
			Name:                             "parse refresh token unknown error",
			Input:                            dto.LoginByRefreshTokenInput{UserID: 1, RefreshToken: "CoolRefreshToken"},
			WantedOutput:                     dto.LoginByRefreshTokenOutput{},
			WantedError:                      test.UnknownTestErr,
			GetUserByIDWantedID:              1,
			GetUserByIDReturningUser:         &domain.User{ID: 1, Username: "Nimk"},
			GetUserByIDReturningError:        nil,
			ExpectRefreshTokenExists:         true,
			RefreshTokenExistsWantedToken:    "CoolRefreshToken",
			RefreshTokenExistsReturningBool:  true,
			RefreshTokenExistsReturningError: nil,
			ExpectParseRefreshToken:          true,
			ParseRefreshTokenWantedToken:     "CoolRefreshToken",
			ParseRefreshTokenReturningID:     0,
			ParseRefreshTokenReturningError:  test.UnknownTestErr,
		},
		{
			Name:                                   "generate token pair unknown error",
			Input:                                  dto.LoginByRefreshTokenInput{UserID: 1, RefreshToken: "CoolRefreshToken"},
			WantedOutput:                           dto.LoginByRefreshTokenOutput{},
			WantedError:                            test.UnknownTestErr,
			GetUserByIDWantedID:                    1,
			GetUserByIDReturningUser:               &domain.User{ID: 1, Username: "Nimk"},
			GetUserByIDReturningError:              nil,
			ExpectRefreshTokenExists:               true,
			RefreshTokenExistsWantedToken:          "CoolRefreshToken",
			RefreshTokenExistsReturningBool:        true,
			RefreshTokenExistsReturningError:       nil,
			ExpectParseRefreshToken:                true,
			ParseRefreshTokenWantedToken:           "CoolRefreshToken",
			ParseRefreshTokenReturningID:           1,
			ParseRefreshTokenReturningError:        nil,
			ExpectGenerateTokenPair:                true,
			GenerateTokenPairWantedID:              1,
			GenerateTokenPairWantedUsername:        "Nimk",
			GenerateTokenPairReturningAccessToken:  "",
			GenerateTokenPairReturningRefreshToken: "",
			GenerateTokenPairReturningError:        test.UnknownTestErr,
		},
		{
			Name:                                   "replace refresh token error not found",
			Input:                                  dto.LoginByRefreshTokenInput{UserID: 1, RefreshToken: "CoolRefreshToken"},
			WantedOutput:                           dto.LoginByRefreshTokenOutput{},
			WantedError:                            domain.ErrNotFound,
			GetUserByIDWantedID:                    1,
			GetUserByIDReturningUser:               &domain.User{ID: 1, Username: "Nimk"},
			GetUserByIDReturningError:              nil,
			ExpectRefreshTokenExists:               true,
			RefreshTokenExistsWantedToken:          "CoolRefreshToken",
			RefreshTokenExistsReturningBool:        true,
			RefreshTokenExistsReturningError:       nil,
			ExpectParseRefreshToken:                true,
			ParseRefreshTokenWantedToken:           "CoolRefreshToken",
			ParseRefreshTokenReturningID:           1,
			ParseRefreshTokenReturningError:        nil,
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
			Input:                                  dto.LoginByRefreshTokenInput{UserID: 1, RefreshToken: "CoolRefreshToken"},
			WantedOutput:                           dto.LoginByRefreshTokenOutput{},
			WantedError:                            test.UnknownTestErr,
			GetUserByIDWantedID:                    1,
			GetUserByIDReturningUser:               &domain.User{ID: 1, Username: "Nimk"},
			GetUserByIDReturningError:              nil,
			ExpectRefreshTokenExists:               true,
			RefreshTokenExistsWantedToken:          "CoolRefreshToken",
			RefreshTokenExistsReturningBool:        true,
			RefreshTokenExistsReturningError:       nil,
			ExpectParseRefreshToken:                true,
			ParseRefreshTokenWantedToken:           "CoolRefreshToken",
			ParseRefreshTokenReturningID:           1,
			ParseRefreshTokenReturningError:        nil,
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

			if tt.ExpectRefreshTokenExists {
				rm.EXPECT().RefreshTokenExists(gomock.Any()).Return(tt.RefreshTokenExistsReturningBool, tt.RefreshTokenExistsReturningError)
			}

			if tt.ExpectParseRefreshToken {
				tim.EXPECT().ParseRefreshToken(tt.ParseRefreshTokenWantedToken).Return(tt.ParseRefreshTokenReturningID, tt.ParseRefreshTokenReturningError)
			}

			if tt.ExpectGenerateTokenPair {
				tim.EXPECT().GenerateTokenPair(tt.GenerateTokenPairWantedID, tt.GenerateTokenPairWantedUsername).Return(tt.GenerateTokenPairReturningAccessToken, tt.GenerateTokenPairReturningRefreshToken, tt.GenerateTokenPairReturningError)
			}

			if tt.ExpectReplaceRefreshToken {
				rm.EXPECT().ReplaceRefreshToken(gomock.Cond(test.CheckRefreshTokenUserID(tt.ReplaceRefreshTokenWantedToken))).Return(tt.ReplaceRefreshTokenReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.LoginByRefreshToken(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
