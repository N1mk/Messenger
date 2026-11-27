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

func TestAddChUserToChat(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.AddUserToChatInput
		WantedOutput dto.AddUserToChatOutput
		WantedError  error

		GetUserChatInvitesWantetID         int
		GetUserChatInvitesReturningInvites []*domain.ChatInvite
		GetUserChatInvitesReturningError   error

		ExpectAddUserToChat         bool
		AddUserToChatWantedUserID   int
		AddUserToChatWantedChatID   int
		AddUserToChatReturningError error
	}{
		{
			Name:                               "ok",
			Input:                              dto.AddUserToChatInput{UserID: 1, ChatID: 1},
			WantedOutput:                       dto.AddUserToChatOutput{},
			WantedError:                        nil,
			GetUserChatInvitesWantetID:         1,
			GetUserChatInvitesReturningInvites: []*domain.ChatInvite{&domain.ChatInvite{ChatID: 1}},
			GetUserChatInvitesReturningError:   nil,
			ExpectAddUserToChat:                true,
			AddUserToChatWantedUserID:          1,
			AddUserToChatWantedChatID:          1,
			AddUserToChatReturningError:        nil,
		},
		{
			Name:                               "error user is not invited",
			Input:                              dto.AddUserToChatInput{UserID: 1, ChatID: 1},
			WantedOutput:                       dto.AddUserToChatOutput{},
			WantedError:                        domain.ErrUserIsNotInvited,
			GetUserChatInvitesWantetID:         1,
			GetUserChatInvitesReturningInvites: []*domain.ChatInvite{&domain.ChatInvite{ChatID: 2}},
			GetUserChatInvitesReturningError:   nil,
		},
		{
			Name:                               "get user chat invites unknown error",
			Input:                              dto.AddUserToChatInput{UserID: 1, ChatID: 1},
			WantedOutput:                       dto.AddUserToChatOutput{},
			WantedError:                        test.UnknownTestErr,
			GetUserChatInvitesWantetID:         1,
			GetUserChatInvitesReturningInvites: []*domain.ChatInvite{},
			GetUserChatInvitesReturningError:   test.UnknownTestErr,
		},
		{
			Name:                               "add user to chat unknown error",
			Input:                              dto.AddUserToChatInput{UserID: 1, ChatID: 1},
			WantedOutput:                       dto.AddUserToChatOutput{},
			WantedError:                        test.UnknownTestErr,
			GetUserChatInvitesWantetID:         1,
			GetUserChatInvitesReturningInvites: []*domain.ChatInvite{&domain.ChatInvite{ChatID: 1}},
			GetUserChatInvitesReturningError:   nil,
			ExpectAddUserToChat:                true,
			AddUserToChatWantedUserID:          1,
			AddUserToChatWantedChatID:          1,
			AddUserToChatReturningError:        test.UnknownTestErr,
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

			rm.EXPECT().GetUserChatInvites(tt.GetUserChatInvitesWantetID).Return(tt.GetUserChatInvitesReturningInvites, tt.GetUserChatInvitesReturningError)

			if tt.ExpectAddUserToChat {
				rm.EXPECT().AddUserToChat(tt.AddUserToChatWantedUserID, tt.AddUserToChatWantedChatID).Return(tt.AddUserToChatReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, mbm, log, tim)

			output, err := u.AddUserToChat(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
