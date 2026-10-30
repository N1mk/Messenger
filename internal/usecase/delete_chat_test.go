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

func TestDeleteChat(t *testing.T) {
	tests := []struct {
		Name          string
		Users         []domain.User
		Chats         []domain.Chat
		Input         dto.DeleteChatInput
		ExpectedError error
		ExpectedChats []domain.ChatWithoutSlices
	}{
		{
			Name:          "Ok",
			Users:         []domain.User{domain.User{ID: 0, Username: "Nimk"}},
			Chats:         []domain.Chat{domain.Chat{ID: 0, Name: "The chat"}},
			Input:         dto.DeleteChatInput{UserID: 0, ChatID: 0},
			ExpectedError: nil,
			ExpectedChats: []domain.ChatWithoutSlices{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Chats[0].Users = test.Users
			test.Users[0].Chats = test.Chats

			fp := postgres.NewFakePostgres()
			fp.Users = test.Users
			fp.Chats = test.Chats

			uc := NewUsecase(fp)

			_, err := uc.DeleteChat(context.TODO(), test.Input)

			if !errors.Is(err, test.ExpectedError) {
				t.Errorf("Expected error: %s | Got: %s", test.ExpectedError.Error(), err.Error())
			}

			if len(fp.Chats) != len(test.ExpectedChats) {
				t.Errorf("Expected Chats: %v | Got: %v", test.ExpectedChats, fp.Chats)
			}

			var chats []domain.ChatWithoutSlices
			for _, chat := range fp.Chats {
				chats = append(chats, domain.ChatWithoutSlices{ID: chat.ID, Name: chat.Name})
			}

			if !sliceoperations.CompareSlices(chats, test.ExpectedChats) {
				t.Errorf("Expected Chats: %v | Got: %v", test.ExpectedChats, chats)
			}
		})
	}
}
