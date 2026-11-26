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

func TestListenToMessages(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.ListenToMessagesInput
		WantedOutput dto.ListenToMessagesOutput
		WantedError  error

		GetUserByIDWantedID       int
		GetUserByIDReturningUser  *domain.User
		GetUserByIDReturningError error

		ExpectMakeUserListenChats         bool
		MakeUserListenChatsWantedRequest  domain.Request
		MakeUserListenChatsWantedUser     *domain.User
		MakeUserListenChatsWantedChats    []*domain.Chat
		MakeUserListenChatsReturningError error
	}{
		{
			Name:                              "ok",
			Input:                             dto.ListenToMessagesInput{UserID: 1, Request: test.TestRequest{Wanted: "wanted data"}},
			WantedOutput:                      dto.ListenToMessagesOutput{},
			WantedError:                       nil,
			GetUserByIDWantedID:               1,
			GetUserByIDReturningUser:          &domain.User{ID: 1, Username: "Nimk", Chats: []domain.Chat{domain.Chat{ID: 1, Name: "First chat"}, domain.Chat{ID: 2, Name: "Second chat"}}},
			GetUserByIDReturningError:         nil,
			ExpectMakeUserListenChats:         true,
			MakeUserListenChatsWantedRequest:  test.TestRequest{Wanted: "wanted data"},
			MakeUserListenChatsWantedUser:     &domain.User{ID: 1, Username: "Nimk", Chats: []domain.Chat{domain.Chat{ID: 1, Name: "First chat"}, domain.Chat{ID: 2, Name: "Second chat"}}},
			MakeUserListenChatsWantedChats:    []*domain.Chat{&domain.Chat{ID: 1, Name: "First chat"}, &domain.Chat{ID: 2, Name: "Second chat"}},
			MakeUserListenChatsReturningError: nil,
		},
		{
			Name:                      "get user by id error not found",
			Input:                     dto.ListenToMessagesInput{UserID: 1, Request: test.TestRequest{Wanted: "wanted data"}},
			WantedOutput:              dto.ListenToMessagesOutput{},
			WantedError:               domain.ErrNotFound,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{ID: 1, Username: "Nimk", Chats: []domain.Chat{domain.Chat{ID: 1, Name: "First chat"}, domain.Chat{ID: 2, Name: "Second chat"}}},
			GetUserByIDReturningError: domain.ErrNotFound,
		},
		{
			Name:                      "get user by id unknown error",
			Input:                     dto.ListenToMessagesInput{UserID: 1, Request: test.TestRequest{Wanted: "wanted data"}},
			WantedOutput:              dto.ListenToMessagesOutput{},
			WantedError:               test.UnknownTestErr,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{ID: 1, Username: "Nimk", Chats: []domain.Chat{domain.Chat{ID: 1, Name: "First chat"}, domain.Chat{ID: 2, Name: "Second chat"}}},
			GetUserByIDReturningError: test.UnknownTestErr,
		},
		{
			Name:                              "make user listen chats unknown error",
			Input:                             dto.ListenToMessagesInput{UserID: 1, Request: test.TestRequest{Wanted: "wanted data"}},
			WantedOutput:                      dto.ListenToMessagesOutput{},
			WantedError:                       test.UnknownTestErr,
			GetUserByIDWantedID:               1,
			GetUserByIDReturningUser:          &domain.User{ID: 1, Username: "Nimk", Chats: []domain.Chat{domain.Chat{ID: 1, Name: "First chat"}, domain.Chat{ID: 2, Name: "Second chat"}}},
			GetUserByIDReturningError:         nil,
			ExpectMakeUserListenChats:         true,
			MakeUserListenChatsWantedRequest:  test.TestRequest{Wanted: "wanted data"},
			MakeUserListenChatsWantedUser:     &domain.User{ID: 1, Username: "Nimk", Chats: []domain.Chat{domain.Chat{ID: 1, Name: "First chat"}, domain.Chat{ID: 2, Name: "Second chat"}}},
			MakeUserListenChatsWantedChats:    []*domain.Chat{&domain.Chat{ID: 1, Name: "First chat"}, &domain.Chat{ID: 2, Name: "Second chat"}},
			MakeUserListenChatsReturningError: test.UnknownTestErr,
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

			rm.EXPECT().GetUserByID(tt.GetUserByIDWantedID).Return(tt.GetUserByIDReturningUser, tt.GetUserByIDReturningError)

			if tt.ExpectMakeUserListenChats {
				mbm.EXPECT().MakeUserListenChats(context.TODO(), tt.MakeUserListenChatsWantedRequest, tt.MakeUserListenChatsWantedUser, tt.MakeUserListenChatsWantedChats).Return(tt.MakeUserListenChatsReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, mbm, log, tim)

			output, err := u.ListenToMessages(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
