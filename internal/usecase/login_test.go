package usecase

import (
	"errors"
	postgres "messenger/internal/adapter/postgres/fake"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"testing"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		Name          string
		Users         []domain.User
		Input         dto.LoginInput
		ExpectedError error
	}{
		{
			Name:          "Ok",
			Users:         []domain.User{domain.User{Username: "Nimk", Email: "nimk@mail.com", Password: "1234"}},
			Input:         dto.LoginInput{UserID: 0, Password: "1234"},
			ExpectedError: nil,
		},
		{
			Name:          "Not found",
			Users:         []domain.User{},
			Input:         dto.LoginInput{UserID: 0, Password: "1234"},
			ExpectedError: domain.ErrNotFound,
		},
		{
			Name:          "Incorrect password",
			Users:         []domain.User{domain.User{Username: "Nimk", Email: "nimk@mail.com", Password: "1234"}},
			Input:         dto.LoginInput{UserID: 0, Password: "12345"},
			ExpectedError: domain.ErrIncorrectPassword,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			fp := postgres.NewFakePostgres()
			fp.Users = test.Users

			uc := NewUsecase(fp)

			_, err := uc.Login(test.Input)

			if !errors.Is(err, test.ExpectedError) {
				t.Errorf("Expected error: %s | Got: %s", test.ExpectedError.Error(), err.Error())
			}
		})
	}
}
