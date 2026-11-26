package messagebroker

import (
	"context"
	"fmt"
	apphttp "messenger/internal/controller/http"
	"messenger/internal/domain"
	"messenger/pkg/logger"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var webSocketupgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type DefaultMessageBroker struct {
	r *redis.Client
	l *logger.Logger

	// userConnectionCancels represents all active websocket conntection with users.
	// Key is a user ID and value is a cancel func of a context of a gorutine which runs the connection.
	userConnectionCancels map[int]context.CancelFunc
}

func NewDefaultMessageBroker(redisAddr string, redisPassword string, redisDB int, logger *logger.Logger) *DefaultMessageBroker {
	return &DefaultMessageBroker{r: redis.NewClient(&redis.Options{Addr: redisAddr, Password: redisPassword, DB: redisDB}), l: logger}
}

func (dmb *DefaultMessageBroker) SendMessage(ctx context.Context, chatID int, text string) error {
	return dmb.r.Publish(ctx, fmt.Sprintf("messages:%d", chatID), text).Err()
}

func (dmb *DefaultMessageBroker) MakeUserListenChats(ctx context.Context, request domain.Request, userID int, chats []*domain.Chat) error {
	close, exists := dmb.userConnectionCancels[userID]
	if exists {
		close()
	}

	httpRequest, ok := request.(apphttp.HTTPRequest)
	if !ok {
		return fmt.Errorf("request type assertion error: request is not of type HTTPRequest")
	}

	conn, err := webSocketupgrader.Upgrade(httpRequest.ResponseWriter, httpRequest.Request, nil)
	if err != nil {
		return fmt.Errorf("upgrading to websocket error: %w", err)
	}

	channels := make([]string, len(chats))
	for i, chat := range chats {
		channels[i] = fmt.Sprintf("messages:%d", chat.ID)
	}

	ps := dmb.r.Subscribe(ctx, channels...)

	messageChan := ps.Channel()

	connectionContext, connectionCancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case message := <-messageChan:
				if err := conn.WriteJSON(struct {
					ChatID string `json:"chat_id"`
					Text   string `json:"text"`
				}{ChatID: strings.TrimPrefix(message.Pattern, "messages:"), Text: message.Payload}); err != nil {
					dmb.l.LogError(fmt.Sprintf("user %d websocket conntection gorutine error: message writing error: %s", userID, err.Error()))
				}
			case <-connectionContext.Done():
				dmb.l.LogInfo(fmt.Sprintf("user %d websocket conntection gorutine error: closing connection", userID))
			}
		}
	}()

	dmb.userConnectionCancels[userID] = connectionCancel

	return nil
}
