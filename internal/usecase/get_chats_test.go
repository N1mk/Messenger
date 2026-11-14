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

func TestGetChats(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.GetChatsInput
		WantedOutput dto.GetChatsOutput
		WantedError  error

		GetUserByIDWantedID       int
		GetUserByIDReturningUser  *domain.User
		GetUserByIDReturningError error
	}{
		{
			Name:                      "ok without chats",
			Input:                     dto.GetChatsInput{UserID: 1},
			WantedOutput:              dto.GetChatsOutput{Chats: []*domain.Chat{}},
			WantedError:               nil,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{Chats: []domain.Chat{}},
			GetUserByIDReturningError: nil,
		},
		{
			Name:                      "ok with chats",
			Input:                     dto.GetChatsInput{UserID: 1},
			WantedOutput:              dto.GetChatsOutput{Chats: []*domain.Chat{&domain.Chat{ID: 1, Name: "First chat"}, &domain.Chat{ID: 2, Name: "Second chat"}}},
			WantedError:               nil,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{Chats: []domain.Chat{domain.Chat{ID: 1, Name: "First chat"}, domain.Chat{ID: 2, Name: "Second chat"}}},
			GetUserByIDReturningError: nil,
		},
		{
			Name:                      "err not found",
			Input:                     dto.GetChatsInput{UserID: 1},
			WantedOutput:              dto.GetChatsOutput{Chats: []*domain.Chat{}},
			WantedError:               domain.ErrNotFound,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{Chats: []domain.Chat{}},
			GetUserByIDReturningError: domain.ErrNotFound,
		},
		{
			Name:                      "unknown error",
			Input:                     dto.GetChatsInput{UserID: 1},
			WantedOutput:              dto.GetChatsOutput{Chats: []*domain.Chat{}},
			WantedError:               test.UnknownTestErr,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{Chats: []domain.Chat{}},
			GetUserByIDReturningError: test.UnknownTestErr,
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

			rm.EXPECT().GetUserByID(tt.GetUserByIDWantedID).Return(tt.GetUserByIDReturningUser, tt.GetUserByIDReturningError)

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.GetChats(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if !sliceoperations.CompareSlicesOfPointers(output.Chats, tt.WantedOutput.Chats) {
				subT.Errorf("wanted chats: %v, got chats: %v", tt.WantedOutput.Chats, output.Chats)
			}
		})
	}
}
