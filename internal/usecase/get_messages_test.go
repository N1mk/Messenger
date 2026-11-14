package usecase

import (
	"context"
	"errors"
	"messenger/internal/domain"
	"messenger/internal/dto"
	"messenger/internal/usecase/mocks"
	"messenger/internal/usecase/test"
	"messenger/pkg/logger"
	sliceoperations "messenger/pkg/slice_operations"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGetMessages(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.GetMessagesInput
		WantedOutput dto.GetMessagesOutput
		WantedError  error

		GetChatByIDWantedID       int
		GetChatByIDReturningChat  *domain.Chat
		GetChatByIDReturningError error

		ExpectGetChatMessages            bool
		GetChatMessagesWantedID          int
		GetChatMessagesWantedLastID      int
		GetChatMessagesWantedPageSize    int
		GetChatMessagesReturningMessages []*domain.Message
		GetChatMessagesReturningLastID   int
		GetChatMessagesRetuningError     error
	}{
		{
			Name:                             "ok without messages",
			Input:                            dto.GetMessagesInput{UserID: 1, ChatID: 1, LastMessageID: 0, PageSize: 0},
			WantedOutput:                     dto.GetMessagesOutput{Messages: []*domain.Message{}, LastMessageID: 0},
			WantedError:                      nil,
			GetChatByIDWantedID:              1,
			GetChatByIDReturningChat:         &domain.Chat{Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError:        nil,
			ExpectGetChatMessages:            true,
			GetChatMessagesWantedID:          1,
			GetChatMessagesWantedLastID:      0,
			GetChatMessagesWantedPageSize:    0,
			GetChatMessagesReturningMessages: []*domain.Message{},
			GetChatMessagesReturningLastID:   0,
			GetChatMessagesRetuningError:     nil,
		},
		{
			Name:                             "ok with messages",
			Input:                            dto.GetMessagesInput{UserID: 1, ChatID: 1, LastMessageID: 0, PageSize: 0},
			WantedOutput:                     dto.GetMessagesOutput{Messages: []*domain.Message{&domain.Message{ID: 1, Text: "Hello world!"}}, LastMessageID: 1},
			WantedError:                      nil,
			GetChatByIDWantedID:              1,
			GetChatByIDReturningChat:         &domain.Chat{Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError:        nil,
			ExpectGetChatMessages:            true,
			GetChatMessagesWantedID:          1,
			GetChatMessagesWantedLastID:      0,
			GetChatMessagesWantedPageSize:    0,
			GetChatMessagesReturningMessages: []*domain.Message{&domain.Message{ID: 1, Text: "Hello world!"}},
			GetChatMessagesReturningLastID:   1,
			GetChatMessagesRetuningError:     nil,
		},
		{
			Name:                      "get chat by id error not found",
			Input:                     dto.GetMessagesInput{UserID: 1, ChatID: 1, LastMessageID: 0, PageSize: 0},
			WantedOutput:              dto.GetMessagesOutput{Messages: []*domain.Message{}, LastMessageID: 0},
			WantedError:               domain.ErrNotFound,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{Users: []domain.User{}},
			GetChatByIDReturningError: domain.ErrNotFound,
		},
		{
			Name:                      "get chat by id unknown error",
			Input:                     dto.GetMessagesInput{UserID: 1, ChatID: 1, LastMessageID: 0, PageSize: 0},
			WantedOutput:              dto.GetMessagesOutput{Messages: []*domain.Message{}, LastMessageID: 0},
			WantedError:               test.UnknownTestErr,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{Users: []domain.User{}},
			GetChatByIDReturningError: test.UnknownTestErr,
		},
		{
			Name:                      "error user is not in chat",
			Input:                     dto.GetMessagesInput{UserID: 1, ChatID: 1, LastMessageID: 0, PageSize: 0},
			WantedOutput:              dto.GetMessagesOutput{Messages: []*domain.Message{}, LastMessageID: 0},
			WantedError:               domain.ErrUserIsNotInChat,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{Users: []domain.User{}},
			GetChatByIDReturningError: nil,
		},
		{
			Name:                             "get chat messages error not found",
			Input:                            dto.GetMessagesInput{UserID: 1, ChatID: 1, LastMessageID: 0, PageSize: 0},
			WantedOutput:                     dto.GetMessagesOutput{Messages: []*domain.Message{}, LastMessageID: 0},
			WantedError:                      domain.ErrNotFound,
			GetChatByIDWantedID:              1,
			GetChatByIDReturningChat:         &domain.Chat{Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError:        nil,
			ExpectGetChatMessages:            true,
			GetChatMessagesWantedID:          1,
			GetChatMessagesWantedLastID:      0,
			GetChatMessagesWantedPageSize:    0,
			GetChatMessagesReturningMessages: []*domain.Message{},
			GetChatMessagesReturningLastID:   1,
			GetChatMessagesRetuningError:     domain.ErrNotFound,
		},
		{
			Name:                             "get chat messages unknown error",
			Input:                            dto.GetMessagesInput{UserID: 1, ChatID: 1, LastMessageID: 0, PageSize: 0},
			WantedOutput:                     dto.GetMessagesOutput{Messages: []*domain.Message{}, LastMessageID: 0},
			WantedError:                      test.UnknownTestErr,
			GetChatByIDWantedID:              1,
			GetChatByIDReturningChat:         &domain.Chat{Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError:        nil,
			ExpectGetChatMessages:            true,
			GetChatMessagesWantedID:          1,
			GetChatMessagesWantedLastID:      0,
			GetChatMessagesWantedPageSize:    0,
			GetChatMessagesReturningMessages: []*domain.Message{},
			GetChatMessagesReturningLastID:   1,
			GetChatMessagesRetuningError:     test.UnknownTestErr,
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

			rm.EXPECT().GetChatByID(tt.GetChatByIDWantedID).Return(tt.GetChatByIDReturningChat, tt.GetChatByIDReturningError)

			if tt.ExpectGetChatMessages {
				rm.EXPECT().GetChatMessages(tt.GetChatMessagesWantedID, tt.GetChatMessagesWantedLastID, tt.GetChatMessagesWantedPageSize).Return(tt.GetChatMessagesReturningMessages, tt.GetChatMessagesReturningLastID, tt.GetChatMessagesRetuningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.GetMessages(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if !sliceoperations.CompareSlicesOfPointers(output.Messages, tt.WantedOutput.Messages) {
				subT.Errorf("wanted messages: %v, got messages: %v", tt.WantedOutput.Messages, output.Messages)
			}

			if output.LastMessageID != tt.WantedOutput.LastMessageID {
				subT.Errorf("wanted lastMessageID: %d, got lastMessageID: %d", tt.WantedOutput.LastMessageID, output.LastMessageID)
			}
		})
	}
}
