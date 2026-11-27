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

func TestSendMessage(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.SendMessageInput
		WantedOutput dto.SendMessageOutput
		WantedError  error

		GetChatByIDWantedID       int
		GetChatByIDReturningChat  *domain.Chat
		GetChatByIDReturningError error

		ExpectAddMessage           bool
		AddMessageWantedMessage    domain.Message
		AddMessageReturningMessage *domain.Message
		AddMessageReturningError   error

		ExpectSendMessage         bool
		SendMessageWantedChatID   int
		SendMessageWantedText     string
		SendMessageReturningError error
	}{
		{
			Name:                       "ok",
			Input:                      dto.SendMessageInput{UserID: 1, ChatID: 1, Text: "Hello world!"},
			WantedOutput:               dto.SendMessageOutput{ID: 1},
			WantedError:                nil,
			GetChatByIDWantedID:        1,
			GetChatByIDReturningChat:   &domain.Chat{Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError:  nil,
			ExpectAddMessage:           true,
			AddMessageWantedMessage:    domain.Message{ChatID: 1, UserID: 1, Text: "Hello world!"},
			AddMessageReturningMessage: &domain.Message{ID: 1, ChatID: 1, UserID: 1, Text: "Hello world!"},
			AddMessageReturningError:   nil,
			ExpectSendMessage:          true,
			SendMessageWantedChatID:    1,
			SendMessageWantedText:      "Hello world!",
			SendMessageReturningError:  nil,
		},
		{
			Name:                      "get chat by id error not found",
			Input:                     dto.SendMessageInput{UserID: 1, ChatID: 1, Text: "Hello world!"},
			WantedOutput:              dto.SendMessageOutput{},
			WantedError:               domain.ErrNotFound,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{Users: []domain.User{}},
			GetChatByIDReturningError: domain.ErrNotFound,
		},
		{
			Name:                      "get chat by id unknown error",
			Input:                     dto.SendMessageInput{UserID: 1, ChatID: 1, Text: "Hello world!"},
			WantedOutput:              dto.SendMessageOutput{},
			WantedError:               test.UnknownTestErr,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{Users: []domain.User{}},
			GetChatByIDReturningError: test.UnknownTestErr,
		},
		{
			Name:                      "error user is not in chat",
			Input:                     dto.SendMessageInput{UserID: 1, ChatID: 1, Text: "Hello world!"},
			WantedOutput:              dto.SendMessageOutput{},
			WantedError:               domain.ErrUserIsNotInChat,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{Users: []domain.User{}},
			GetChatByIDReturningError: nil,
		},
		{
			Name:                       "add message unknown error",
			Input:                      dto.SendMessageInput{UserID: 1, ChatID: 1, Text: "Hello world!"},
			WantedOutput:               dto.SendMessageOutput{},
			WantedError:                test.UnknownTestErr,
			GetChatByIDWantedID:        1,
			GetChatByIDReturningChat:   &domain.Chat{Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError:  nil,
			ExpectAddMessage:           true,
			AddMessageWantedMessage:    domain.Message{ChatID: 1, UserID: 1, Text: "Hello world!"},
			AddMessageReturningMessage: nil,
			AddMessageReturningError:   test.UnknownTestErr,
		},
		{
			Name:                       "send message unknown error",
			Input:                      dto.SendMessageInput{UserID: 1, ChatID: 1, Text: "Hello world!"},
			WantedOutput:               dto.SendMessageOutput{},
			WantedError:                test.UnknownTestErr,
			GetChatByIDWantedID:        1,
			GetChatByIDReturningChat:   &domain.Chat{Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError:  nil,
			ExpectAddMessage:           true,
			AddMessageWantedMessage:    domain.Message{ChatID: 1, UserID: 1, Text: "Hello world!"},
			AddMessageReturningMessage: &domain.Message{ID: 1, ChatID: 1, UserID: 1, Text: "Hello world!"},
			AddMessageReturningError:   nil,
			ExpectSendMessage:          true,
			SendMessageWantedChatID:    1,
			SendMessageWantedText:      "Hello world!",
			SendMessageReturningError:  test.UnknownTestErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(subT *testing.T) {
			ctrl := gomock.NewController(subT)
			rm := mocks.NewMockRepo(ctrl)
			mbm := mocks.NewMockMessageBroker(ctrl)
			tim := mocks.NewMockTokenIssuer(ctrl)

			rm.EXPECT().Begin(context.TODO())

			if tt.WantedError == nil {
				rm.EXPECT().Commit()
			} else {
				rm.EXPECT().Rollback()
			}

			rm.EXPECT().GetChatByID(tt.GetChatByIDWantedID).Return(tt.GetChatByIDReturningChat, tt.GetChatByIDReturningError)

			if tt.ExpectAddMessage {
				rm.EXPECT().AddMessage(gomock.Cond(test.CheckMessageTextAndAndChatIDAndUserID(tt.AddMessageWantedMessage))).Return(tt.AddMessageReturningMessage, tt.AddMessageReturningError)
			}

			if tt.ExpectSendMessage {
				mbm.EXPECT().SendMessage(context.TODO(), tt.SendMessageWantedChatID, tt.SendMessageWantedText).Return(tt.SendMessageReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, mbm, log, tim)

			output, err := u.SendMessage(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
