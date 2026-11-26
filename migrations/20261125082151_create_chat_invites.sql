-- +goose Up
CREATE TABLE chat_invites (
    chat_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL
);

-- +goose Down
DROP TABLE chat_invites;
