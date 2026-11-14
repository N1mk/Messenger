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

func TestUpdateChat(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.UpdateChatInput
		WantedOutput dto.UpdateChatOutput
		WantedError  error

		GetChatByIDWantedID       int
		GetChatByIDReturningChat  *domain.Chat
		GetChatByIDReturningError error

		ExpectGetUserByName         bool
		GetUserByNameNames          []string
		GetUserByNameReturningError error

		ExpectUpdateChat         bool
		UpdateChatWantedChat     domain.Chat
		UpdateChatReturningError error
	}{
		{
			Name:                      "ok without users",
			Input:                     dto.UpdateChatInput{UserID: 1, ChatID: 1, OwnerID: 1, Usernames: []string{}, Name: "New chat name"},
			WantedOutput:              dto.UpdateChatOutput{},
			WantedError:               nil,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{OwnerID: 1},
			GetChatByIDReturningError: nil,
			ExpectUpdateChat:          true,
			UpdateChatWantedChat:      domain.Chat{ID: 1, OwnerID: 1, Name: "New chat name", Users: []domain.User{}},
			UpdateChatReturningError:  nil,
		},
		{
			Name:                        "ok with users",
			Input:                       dto.UpdateChatInput{UserID: 1, ChatID: 1, OwnerID: 1, Usernames: []string{"Nimk", "Second user"}, Name: "New chat name"},
			WantedOutput:                dto.UpdateChatOutput{},
			WantedError:                 nil,
			GetChatByIDWantedID:         1,
			GetChatByIDReturningChat:    &domain.Chat{OwnerID: 1},
			GetChatByIDReturningError:   nil,
			ExpectGetUserByName:         true,
			GetUserByNameNames:          []string{"Nimk", "Second user"},
			GetUserByNameReturningError: nil,
			ExpectUpdateChat:            true,
			UpdateChatWantedChat:        domain.Chat{ID: 1, OwnerID: 1, Name: "New chat name", Users: []domain.User{domain.User{ID: 1, Username: "Nimk"}, domain.User{ID: 2, Username: "Second user"}}},
			UpdateChatReturningError:    nil,
		},
		{
			Name:                      "get chat by id error not found",
			Input:                     dto.UpdateChatInput{UserID: 1, ChatID: 1, OwnerID: 1, Usernames: []string{}, Name: "New chat name"},
			WantedOutput:              dto.UpdateChatOutput{},
			WantedError:               domain.ErrNotFound,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{OwnerID: 1},
			GetChatByIDReturningError: domain.ErrNotFound,
		},
		{
			Name:                      "get chat by id unknown error",
			Input:                     dto.UpdateChatInput{UserID: 1, ChatID: 1, OwnerID: 1, Usernames: []string{}, Name: "New chat name"},
			WantedOutput:              dto.UpdateChatOutput{},
			WantedError:               test.UnknownTestErr,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{OwnerID: 1},
			GetChatByIDReturningError: test.UnknownTestErr,
		},
		{
			Name:                        "get user by name error not found",
			Input:                       dto.UpdateChatInput{UserID: 1, ChatID: 1, OwnerID: 1, Usernames: []string{"Nimk"}, Name: "New chat name"},
			WantedOutput:                dto.UpdateChatOutput{},
			WantedError:                 domain.ErrNotFound,
			GetChatByIDWantedID:         1,
			GetChatByIDReturningChat:    &domain.Chat{OwnerID: 1},
			GetChatByIDReturningError:   nil,
			ExpectGetUserByName:         true,
			GetUserByNameNames:          []string{"Nimk"},
			GetUserByNameReturningError: domain.ErrNotFound,
		},
		{
			Name:                        "get user by name unknown error",
			Input:                       dto.UpdateChatInput{UserID: 1, ChatID: 1, OwnerID: 1, Usernames: []string{"Nimk"}, Name: "New chat name"},
			WantedOutput:                dto.UpdateChatOutput{},
			WantedError:                 test.UnknownTestErr,
			GetChatByIDWantedID:         1,
			GetChatByIDReturningChat:    &domain.Chat{OwnerID: 1},
			GetChatByIDReturningError:   nil,
			ExpectGetUserByName:         true,
			GetUserByNameNames:          []string{"Nimk"},
			GetUserByNameReturningError: test.UnknownTestErr,
		},
		{
			Name:                      "update chat unknown error",
			Input:                     dto.UpdateChatInput{UserID: 1, ChatID: 1, OwnerID: 1, Usernames: []string{}, Name: "New chat name"},
			WantedOutput:              dto.UpdateChatOutput{},
			WantedError:               test.UnknownTestErr,
			GetChatByIDWantedID:       1,
			GetChatByIDReturningChat:  &domain.Chat{OwnerID: 1},
			GetChatByIDReturningError: nil,
			ExpectUpdateChat:          true,
			UpdateChatWantedChat:      domain.Chat{ID: 1, OwnerID: 1, Name: "New chat name", Users: []domain.User{}},
			UpdateChatReturningError:  test.UnknownTestErr,
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

			if tt.ExpectGetUserByName {
				for i := 0; i < len(tt.GetUserByNameNames); i++ {
					rm.EXPECT().GetUserByName(tt.GetUserByNameNames[i]).Return(&domain.User{ID: i + 1, Username: tt.GetUserByNameNames[i]}, tt.GetUserByNameReturningError)
				}
			}

			if tt.ExpectUpdateChat {
				rm.EXPECT().UpdateChat(tt.UpdateChatWantedChat).Return(tt.UpdateChatReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.UpdateChat(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
