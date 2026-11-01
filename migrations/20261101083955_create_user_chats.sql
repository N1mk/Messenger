-- +goose Up
CREATE TABLE user_chats (
    chat_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL
);

-- +goose Down
DROP TABLE user_chats;
