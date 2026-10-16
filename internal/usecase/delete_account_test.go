package usecase

import (
	"errors"
	postgres "messenger/internal/adapter/postgres/fake"
	"messenger/internal/domain"
	"messenger/internal/dto"
	sliceoperations "messenger/pkg/slice_operations"
	"testing"
)

func TestDeleteAccount(t *testing.T) {
	tests := []struct {
		Name          string
		Users         []domain.User
		Input         dto.DeleteAccountInput
		ExpectedError error
		ExpectedUsers []domain.UserWithoutSlices
	}{
		{
			Name:          "Ok",
			Users:         []domain.User{domain.User{ID: 0, Username: "Nimk", Email: "nimk@mail.com", Password: "1234"}},
			Input:         dto.DeleteAccountInput{UserID: 0},
			ExpectedError: nil,
			ExpectedUsers: []domain.UserWithoutSlices{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			fp := postgres.NewFakePostgres()
			fp.Users = test.Users

			uc := NewUsecase(fp)

			_, err := uc.DeleteAccount(test.Input)

			if !errors.Is(err, test.ExpectedError) {
				t.Errorf("Expected error: %s | Got: %s", test.ExpectedError.Error(), err.Error())
			}

			if len(fp.Users) != len(test.ExpectedUsers) {
				t.Errorf("Expected users: %v | Got: %v", test.ExpectedUsers, fp.Users)
			}

			var users []domain.UserWithoutSlices
			for _, user := range fp.Users {
				users = append(users, domain.UserWithoutSlices{ID: user.ID, Username: user.Username, Email: user.Email, Password: user.Password})
			}

			if !sliceoperations.CompareSlices(users, test.ExpectedUsers) {
				t.Errorf("Expected users: %v | Got: %v", test.ExpectedUsers, users)
			}
		})
	}
}
