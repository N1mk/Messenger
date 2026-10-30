package usecase

import (
	"context"
	"errors"
	postgres "messenger/internal/adapter/postgres/fake"
	"messenger/internal/domain"
	"messenger/internal/dto"
	sliceoperations "messenger/pkg/slice_operations"
	"testing"
	"time"
)

func TestGetMessages(t *testing.T) {
	tests := []struct {
		Name                     string
		Users                    []domain.User
		Chats                    []domain.Chat
		Input                    dto.GetMessagesInput
		ExpectedMessagesInOutput []domain.Message
		ExpectedError            error
	}{
		{
			Name:                     "Ok",
			Users:                    []domain.User{domain.User{ID: 0, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("1234")}},
			Chats:                    []domain.Chat{domain.Chat{ID: 0, Name: "Nimk's chat", Messages: []domain.Message{domain.Message{ChatID: 0, UserID: 0, Text: "Hello, world!"}}}},
			Input:                    dto.GetMessagesInput{UserID: 0, ChatID: 0},
			ExpectedMessagesInOutput: []domain.Message{domain.Message{ChatID: 0, UserID: 0, Text: "Hello, world!", CreatedAt: time.Now()}},
			ExpectedError:            nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Chats[0].Users = test.Users
			test.Users[0].Chats = test.Chats

			fp := postgres.NewFakePostgres()
			fp.Users = test.Users
			fp.Chats = test.Chats

			test.ExpectedMessagesInOutput[0].CreatedAt = test.Chats[0].Messages[0].CreatedAt

			uc := NewUsecase(fp)

			output, err := uc.GetMessages(context.TODO(), test.Input)

			if !errors.Is(err, test.ExpectedError) {
				t.Errorf("Expected error: %s | Got: %s", test.ExpectedError.Error(), err.Error())
			}

			messages := make([]domain.Message, 0)
			for _, message := range output.Messages {
				messages = append(messages, *message)
			}

			if !sliceoperations.CompareSlices(messages, test.ExpectedMessagesInOutput) {
				t.Errorf("Expected messages: %v | Got: %v", test.ExpectedMessagesInOutput, messages)
			}
		})
	}
}
