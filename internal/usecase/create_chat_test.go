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

func TestCreateChat(t *testing.T) {
	tests := []struct {
		Name         string
		Input        dto.CreateChatInput
		WantedOutput dto.CreateChatOutput
		WantedError  error

		GetUserByIDWantedID       int
		GetUserByIDReturningUser  *domain.User
		GetUserByIDReturningError error

		ExpectGetUserByName         bool
		GetUserByNameUsers          []*domain.User
		GetUserByNameReturningError error

		ExpectAddChat         bool
		AddChatWantedChat     domain.Chat
		AddChatReturningChat  *domain.Chat
		AddChatReturningError error
	}{
		{
			Name:                      "without usernames ok",
			Input:                     dto.CreateChatInput{UserID: 1, Usernames: []string{}, Name: "Test chat"},
			WantedOutput:              dto.CreateChatOutput{ChatID: 1},
			WantedError:               nil,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")},
			GetUserByIDReturningError: nil,
			ExpectAddChat:             true,
			AddChatWantedChat:         domain.Chat{OwnerID: 1, Name: "Test chat", Users: []domain.User{domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")}}},
			AddChatReturningChat:      &domain.Chat{ID: 1, OwnerID: 1, Name: "Test chat"},
			AddChatReturningError:     nil,
		},
		{
			Name:                        "with usernames ok",
			Input:                       dto.CreateChatInput{UserID: 1, Usernames: []string{"Another user"}, Name: "Test chat"},
			WantedOutput:                dto.CreateChatOutput{ChatID: 1},
			WantedError:                 nil,
			GetUserByIDWantedID:         1,
			GetUserByIDReturningUser:    &domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")},
			GetUserByIDReturningError:   nil,
			ExpectGetUserByName:         true,
			GetUserByNameUsers:          []*domain.User{&domain.User{ID: 2, Username: "Another user", Email: "anotherUser@mail.com", PasswordHash: []byte("123456")}},
			GetUserByNameReturningError: nil,
			ExpectAddChat:               true,
			AddChatWantedChat:           domain.Chat{OwnerID: 1, Name: "Test chat", Users: []domain.User{domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")}, domain.User{ID: 2, Username: "Another user", Email: "anotherUser@mail.com", PasswordHash: []byte("123456")}}},
			AddChatReturningChat:        &domain.Chat{ID: 1, OwnerID: 1, Name: "Test chat"},
			AddChatReturningError:       nil,
		},
		{
			Name:                      "without usernames get user by id error not found",
			Input:                     dto.CreateChatInput{UserID: 1, Usernames: []string{}, Name: "Test chat"},
			WantedOutput:              dto.CreateChatOutput{},
			WantedError:               domain.ErrNotFound,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{},
			GetUserByIDReturningError: domain.ErrNotFound,
		},
		{
			Name:                      "without usernames get user by id unknown error",
			Input:                     dto.CreateChatInput{UserID: 1, Usernames: []string{}, Name: "Test chat"},
			WantedOutput:              dto.CreateChatOutput{},
			WantedError:               test.UnknownTestErr,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{},
			GetUserByIDReturningError: test.UnknownTestErr,
		},
		{
			Name:                        "with usernames get user by name not found",
			Input:                       dto.CreateChatInput{UserID: 1, Usernames: []string{"Another user"}, Name: "Test chat"},
			WantedOutput:                dto.CreateChatOutput{},
			WantedError:                 domain.ErrNotFound,
			GetUserByIDWantedID:         1,
			GetUserByIDReturningUser:    &domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")},
			GetUserByIDReturningError:   nil,
			ExpectGetUserByName:         true,
			GetUserByNameUsers:          []*domain.User{&domain.User{ID: 2, Username: "Another user", Email: "anotherUser@mail.com", PasswordHash: []byte("123456")}},
			GetUserByNameReturningError: domain.ErrNotFound,
		},
		{
			Name:                        "with usernames get user by name unknown error",
			Input:                       dto.CreateChatInput{UserID: 1, Usernames: []string{"Another user"}, Name: "Test chat"},
			WantedOutput:                dto.CreateChatOutput{},
			WantedError:                 test.UnknownTestErr,
			GetUserByIDWantedID:         1,
			GetUserByIDReturningUser:    &domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")},
			GetUserByIDReturningError:   nil,
			ExpectGetUserByName:         true,
			GetUserByNameUsers:          []*domain.User{&domain.User{ID: 2, Username: "Another user", Email: "anotherUser@mail.com", PasswordHash: []byte("123456")}},
			GetUserByNameReturningError: test.UnknownTestErr,
		},
		{
			Name:                      "without usernames add chat unknown error",
			Input:                     dto.CreateChatInput{UserID: 1, Usernames: []string{}, Name: "Test chat"},
			WantedOutput:              dto.CreateChatOutput{},
			WantedError:               test.UnknownTestErr,
			GetUserByIDWantedID:       1,
			GetUserByIDReturningUser:  &domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")},
			GetUserByIDReturningError: nil,
			ExpectAddChat:             true,
			AddChatWantedChat:         domain.Chat{OwnerID: 1, Name: "Test chat", Users: []domain.User{domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")}}},
			AddChatReturningChat:      &domain.Chat{ID: 1, OwnerID: 1, Name: "Test chat"},
			AddChatReturningError:     test.UnknownTestErr,
		},
		{
			Name:                        "with usernames add chat unknown error",
			Input:                       dto.CreateChatInput{UserID: 1, Usernames: []string{"Another user"}, Name: "Test chat"},
			WantedOutput:                dto.CreateChatOutput{},
			WantedError:                 test.UnknownTestErr,
			GetUserByIDWantedID:         1,
			GetUserByIDReturningUser:    &domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")},
			GetUserByIDReturningError:   nil,
			ExpectGetUserByName:         true,
			GetUserByNameUsers:          []*domain.User{&domain.User{ID: 2, Username: "Another user", Email: "anotherUser@mail.com", PasswordHash: []byte("123456")}},
			GetUserByNameReturningError: nil,
			ExpectAddChat:               true,
			AddChatWantedChat:           domain.Chat{OwnerID: 1, Name: "Test chat", Users: []domain.User{domain.User{ID: 1, Username: "Nimk", Email: "nimk@mail.com", PasswordHash: []byte("123456")}, domain.User{ID: 2, Username: "Another user", Email: "anotherUser@mail.com", PasswordHash: []byte("123456")}}},
			AddChatReturningChat:        &domain.Chat{ID: 1, OwnerID: 1, Name: "Test chat"},
			AddChatReturningError:       test.UnknownTestErr,
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

			if tt.ExpectGetUserByName {
				for _, user := range tt.GetUserByNameUsers {
					rm.EXPECT().GetUserByName(user.Username).Return(user, tt.GetUserByNameReturningError)
				}
			}

			if tt.ExpectAddChat {
				rm.EXPECT().AddChat(tt.AddChatWantedChat).Return(tt.AddChatReturningChat, tt.AddChatReturningError)
			}

			log, _ := logger.NewLogger("", "")

			u := NewUsecase(rm, log, tim)

			output, err := u.CreateChat(context.TODO(), tt.Input)

			if !errors.Is(err, tt.WantedError) {
				subT.Errorf("wanted error: %s, got error: %s", tt.WantedError.Error(), err.Error())
			}

			if output != tt.WantedOutput {
				subT.Errorf("wanted output: %v, got output: %v", tt.WantedOutput, output)
			}
		})
	}
}
