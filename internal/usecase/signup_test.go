package usecase

import (
	"errors"
	postgres "messenger/internal/adapter/postgres/fake"
	"messenger/internal/domain"
	"messenger/internal/dto"
	sliceoperations "messenger/pkg/slice_operations"
	"testing"
)

func TestSignup(t *testing.T) {
	tests := []struct {
		Name          string
		Users         []domain.User
		Input         dto.SignupInput
		ExpectedError error
		ExpectedUsers []domain.UserWithoutSlices
	}{
		{
			Name:          "Ok",
			Users:         []domain.User{},
			Input:         dto.SignupInput{Username: "Nimk", Email: "nimk@mail.com", Password: "1234"},
			ExpectedError: nil,
			ExpectedUsers: []domain.UserWithoutSlices{domain.UserWithoutSlices{ID: 0, Username: "Nimk", Email: "nimk@mail.com", Password: "1234"}},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			fp := postgres.NewFakePostgres()
			fp.Users = test.Users

			uc := NewUsecase(fp)

			_, err := uc.Signup(test.Input)

			if !errors.Is(err, test.ExpectedError) {
				t.Errorf("Expected error: %s | Got: %s", test.ExpectedError.Error(), err.Error())
			}

			if len(fp.Users) != len(test.ExpectedUsers) {
				t.Errorf("Expected users: %v | Got: %v", test.ExpectedUsers, fp.Users)
			}

			var users []domain.UserWithoutSlices
			for _, user := range fp.Users {
				users = append(users, domain.UserWithoutSlices{ID: user.ID, Username: user.Username, Email: user.Email, Password: string(user.PasswordHash)})
			}

			if !sliceoperations.CompareSlices(users, test.ExpectedUsers) {
				t.Errorf("Expected users: %v | Got: %v", test.ExpectedUsers, users)
			}
		})
	}
}
