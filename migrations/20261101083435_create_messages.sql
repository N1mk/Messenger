-- +goose Up
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL REFERENCES chats(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    text VARCHAR(800),
    created_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE chats;