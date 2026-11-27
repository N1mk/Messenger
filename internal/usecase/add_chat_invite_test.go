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

func TestAddChatInvite(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.AddChatInviteInput
		WantedOutput dto.AddChatInviteOutput
		WantedError  error

		GetChatByIDWantedID       int
		GetChatByIDReturningChat  *domain.Chat
		GetChatByIDReturningError error

		ExpectAddChatInvite         bool
		AddChatInviteWantedInvite   domain.ChatInvite
		AddChatInviteReturningError error
	}{
		{
			Name:                        "ok",
			Input:                       dto.AddChatInviteInput{UserID: 1, InvitedID: 2, ChatID: 1},
			WantedOutput:                dto.AddChatInviteOutput{},
			WantedError:                 nil,
			GetChatByIDWantedID:         1,
			GetChatByIDReturningChat:    &domain.Chat{Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError:   nil,
			ExpectAddChatInvite:         true,
			AddChatInviteWantedInvite:   domain.ChatInvite{ChatID: 1, UserID: 2},
			AddChatInviteReturningError: nil,
		},
		{
			Name:                      "error user is not in chat",
			Input:                     dto.AddChatInviteInput{UserID: 1, InvitedID: 2, ChatID: 1},
			WantedOutput:              dto.AddChatInviteOutput{},
			WantedError:               domain.ErrUserIsNotInChat,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{Users: []domain.User{domain.User{ID: 3}}},
			GetChatByIDReturningError: nil,
		},
		{
			Name:                      "get chat by id unknown error",
			Input:                     dto.AddChatInviteInput{UserID: 1, InvitedID: 2, ChatID: 1},
			WantedOutput:              dto.AddChatInviteOutput{},
			WantedError:               test.UnknownTestErr,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  nil,
			GetChatByIDReturningError: test.UnknownTestErr,
		},
		{
			Name:                        "add chat invite unknown error",
			Input:                       dto.AddChatInviteInput{UserID: 1, InvitedID: 2, ChatID: 1},
			WantedOutput:                dto.AddChatInviteOutput{},
			WantedError:                 test.UnknownTestErr,
			GetChatByIDWantedID:         1,
			GetChatByIDReturningChat:    &domain.Chat{Users: []domain.User{domain.User{ID: 1}}},
			GetChatByIDReturningError:   nil,
			ExpectAddChatInvite:         true,
			AddChatInviteWantedInvite:   domain.ChatInvite{ChatID: 1, UserID: 2},
			AddChatInviteReturningError: test.UnknownTestErr,
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

			if tt.ExpectAddChatInvite {
				rm.EXPECT().AddChatInvite(tt.AddChatInviteWantedInvite).Return(tt.AddChatInviteReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, mbm, log, tim)

			output, err := u.AddChatInvite(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
