package usecase

import (
	"errors"
	postgres "messenger/internal/adapter/postgres/fake"
	"messenger/internal/domain"
	"messenger/internal/dto"

	sliceoperations "messenger/pkg/slice_operations"
	"testing"
)

func TestCreateChat(t *testing.T) {
	tests := []struct {
		Name          string
		Users         []domain.User
		Chats         []domain.Chat
		Input         dto.CreateChatInput
		ExpectedError error
		ExpectedChats []domain.ChatWithoutSlices
	}{
		{
			Name:          "Ok",
			Users:         []domain.User{domain.User{Username: "Nimk", Email: "nimk@mail.com", Password: "1234"}},
			Chats:         []domain.Chat{},
			Input:         dto.CreateChatInput{CreatorID: 0, Usernames: []string{"Nimk"}, Name: "Nimk's chat"},
			ExpectedError: nil,
			ExpectedChats: []domain.ChatWithoutSlices{domain.ChatWithoutSlices{ID: 0, Name: "Nimk's chat"}},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			fp := postgres.NewFakePostgres()
			fp.Users = test.Users
			fp.Chats = test.Chats

			uc := NewUsecase(fp)

			_, err := uc.CreateChat(test.Input)

			if !errors.Is(err, test.ExpectedError) {
				t.Errorf("Expected error: %s | Got: %s", test.ExpectedError.Error(), err.Error())
			}

			if len(fp.Chats) != len(test.ExpectedChats) {
				t.Errorf("Expected chats: %v | Got: %v", test.ExpectedChats, fp.Chats)
			}

			var chats []domain.ChatWithoutSlices
			for _, chat := range fp.Chats {
				chats = append(chats, domain.ChatWithoutSlices{ID: chat.ID, Name: chat.Name})
			}

			if !sliceoperations.CompareSlices(chats, test.ExpectedChats) {
				t.Errorf("Expected chats: %v | Got: %v", test.ExpectedChats, chats)
			}
		})
	}
}
