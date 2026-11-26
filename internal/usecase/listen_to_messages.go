package usecase

import (
	"context"
	"fmt"
	"messenger/internal/domain"
	"messenger/internal/dto"
)

// ListenToMessages implements the listen_messages endpoint on the busines layer.
// The method recie dto.GetListenMessages and context.Context and returns dto.ListenToMessages and error.
// It establishes the two-way conntection whith the user, which sends notifications on every chat message.
func (u *Usecase) ListenToMessages(ctx context.Context, input dto.ListenToMessagesInput) (output dto.ListenToMessagesOutput, err error) {
	u.r.Begin(ctx)
	defer func() {
		if err != nil {
			u.r.Rollback()
			return
		}
		u.r.Commit()
	}()

	user, err := u.r.GetUserByID(input.UserID)
	if err != nil {
		u.l.LogError(fmt.Sprintf("user getting error: %s", err.Error()))
		return dto.ListenToMessagesOutput{}, fmt.Errorf("user getting error: %w", err)
	}

	chats := make([]*domain.Chat, 0)
	for _, chat := range user.Chats {
		chats = append(chats, &chat)
	}

	if err := u.mb.MakeUserListenChats(ctx, input.Request, input.UserID, chats); err != nil {
		return dto.ListenToMessagesOutput{}, fmt.Errorf("connaction establishing error: %w", err)
	}

	return dto.ListenToMessagesOutput{}, nil
}
