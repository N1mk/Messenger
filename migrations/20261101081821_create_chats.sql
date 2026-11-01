-- +goose Up
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL REFERENCES users(id),
    name VARCHAR(64) NOT NULL
);

-- +goose Down
DROP TABLE chats;
