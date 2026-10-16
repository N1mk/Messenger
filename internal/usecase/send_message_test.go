package usecase

import (
	"errors"
	postgres "messenger/internal/adapter/postgres/fake"
	"messenger/internal/domain"
	"messenger/internal/dto"
	sliceoperations "messenger/pkg/slice_operations"
	"testing"
	"time"
)

func TestSendMessage(t *testing.T) {
	tests := []struct {
		Name             string
		Users            []domain.User
		Chats            []domain.Chat
		Input            dto.SendMessageInput
		ExpectedMessages []domain.Message
		ExpectedError    error
	}{
		{
			Name:             "Ok",
			Users:            []domain.User{domain.User{ID: 0, Username: "Nimk", Email: "nimk@mail.com", Password: "1234"}},
			Chats:            []domain.Chat{domain.Chat{ID: 0, Name: "Nimk's chat", Messages: make([]domain.Message, 0)}},
			Input:            dto.SendMessageInput{UserID: 0, ChatID: 0, Text: "Hello, world!"},
			ExpectedMessages: []domain.Message{domain.Message{ChatID: 0, UserID: 0, Text: "Hello, world!", CreatedAt: time.Now()}},
			ExpectedError:    nil,
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

			_, err := uc.SendMessage(test.Input)

			if !errors.Is(err, test.ExpectedError) {
				t.Errorf("Expected error: %s | Got: %s", test.ExpectedError.Error(), err.Error())
			}

			fp.Chats[0].Messages[0].CreatedAt = test.ExpectedMessages[0].CreatedAt

			if !sliceoperations.CompareSlices(fp.Chats[0].Messages, test.ExpectedMessages) {
				t.Errorf("Expected messages: %v | Got: %v", test.ExpectedMessages, fp.Chats[0].Messages)
			}
		})
	}
}
