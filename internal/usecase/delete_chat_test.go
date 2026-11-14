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

func TestDeleteChat(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.DeleteChatInput
		WantedOutput dto.DeleteChatOutput
		WantedError  error

		GetChatByIDWantedID       int
		GetChatByIDReturningChat  *domain.Chat
		GetChatByIDReturningError error

		ExpectRemoveChat         bool
		RemoveChatWantedID       int
		RemoveChatReturningError error
	}{
		{
			Name:                      "ok",
			Input:                     dto.DeleteChatInput{UserID: 1, ChatID: 1},
			WantedOutput:              dto.DeleteChatOutput{},
			WantedError:               nil,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{ID: 1, Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError: nil,
			ExpectRemoveChat:          true,
			RemoveChatWantedID:        1,
			RemoveChatReturningError:  nil,
		},
		{
			Name:                      "get chat by id error not found",
			Input:                     dto.DeleteChatInput{UserID: 1, ChatID: 1},
			WantedOutput:              dto.DeleteChatOutput{},
			WantedError:               domain.ErrNotFound,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{ID: 1, Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError: domain.ErrNotFound,
		},
		{
			Name:                      "get chat by id unknown error",
			Input:                     dto.DeleteChatInput{UserID: 1, ChatID: 1},
			WantedOutput:              dto.DeleteChatOutput{},
			WantedError:               test.UnknownTestErr,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{ID: 1, Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError: test.UnknownTestErr,
		},
		{
			Name:                      "error user is not in chat",
			Input:                     dto.DeleteChatInput{UserID: 1, ChatID: 1},
			WantedOutput:              dto.DeleteChatOutput{},
			WantedError:               domain.ErrUserIsNotInChat,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{ID: 1, Users: []domain.User{domain.User{ID: 2}}},
			GetChatByIDReturningError: nil,
		},
		{
			Name:                      "remove chat error not found",
			Input:                     dto.DeleteChatInput{UserID: 1, ChatID: 1},
			WantedOutput:              dto.DeleteChatOutput{},
			WantedError:               domain.ErrNotFound,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{ID: 1, Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError: nil,
			ExpectRemoveChat:          true,
			RemoveChatWantedID:        1,
			RemoveChatReturningError:  domain.ErrNotFound,
		},
		{
			Name:                      "remove chat unknown error",
			Input:                     dto.DeleteChatInput{UserID: 1, ChatID: 1},
			WantedOutput:              dto.DeleteChatOutput{},
			WantedError:               test.UnknownTestErr,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{ID: 1, Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError: nil,
			ExpectRemoveChat:          true,
			RemoveChatWantedID:        1,
			RemoveChatReturningError:  test.UnknownTestErr,
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

			if tt.ExpectRemoveChat {
				rm.EXPECT().RemoveChat(tt.RemoveChatWantedID).Return(tt.RemoveChatReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.DeleteChat(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
