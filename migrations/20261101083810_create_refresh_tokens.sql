-- +goose Up
CREATE TABLE refresh_tokens (
    token_hash BYTEA NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id)
);

-- +goose Down
DROP TABLE refresh_tokens;
