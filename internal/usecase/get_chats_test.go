package usecase

import (
	"context"
	"errors"
	postgres "messenger/internal/adapter/postgres/fake"
	"messenger/internal/domain"
	"messenger/internal/dto"

	sliceoperations "messenger/pkg/slice_operations"
	"testing"
)

func TestGetChats(t *testing.T) {
	tests := []struct {
		Name                    string
		Users                   []domain.User
		Chats                   []domain.Chat
		Input                   dto.GetChatsInput
		ExpectedChatsFromOutput []domain.ChatWithoutSlices
		ExpectedError           error
	}{
		{
			Name:                    "Ok",
			Users:                   []domain.User{domain.User{ID: 0, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("1234"), Chats: []domain.Chat{domain.Chat{ID: 0, Name: "Chat 1"}, domain.Chat{ID: 1, Name: "Chat 2"}}}},
			Chats:                   []domain.Chat{domain.Chat{ID: 0, Name: "Chat 1"}, domain.Chat{ID: 1, Name: "Chat 2"}},
			Input:                   dto.GetChatsInput{UserID: 0},
			ExpectedChatsFromOutput: []domain.ChatWithoutSlices{domain.ChatWithoutSlices{ID: 0, Name: "Chat 1"}, domain.ChatWithoutSlices{ID: 1, Name: "Chat 2"}},
			ExpectedError:           nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			fp := postgres.NewFakePostgres()
			fp.Users = test.Users
			fp.Chats = test.Chats

			uc := NewUsecase(fp)

			output, err := uc.GetChats(context.TODO(), test.Input)

			if !errors.Is(err, test.ExpectedError) {
				t.Errorf("Expected error: %s | Got: %s", test.ExpectedError.Error(), err.Error())
			}

			if len(output.Chats) != len(test.ExpectedChatsFromOutput) {
				t.Errorf("Expected chats: %v | Got: %v", test.ExpectedChatsFromOutput, output.Chats)
			}

			var chats []domain.ChatWithoutSlices
			for _, chat := range output.Chats {
				chats = append(chats, domain.ChatWithoutSlices{ID: chat.ID, Name: chat.Name})
			}

			if !sliceoperations.CompareSlices(chats, test.ExpectedChatsFromOutput) {
				t.Errorf("Expected chats: %v | Got: %v", test.ExpectedChatsFromOutput, chats)
			}
		})
	}
}
