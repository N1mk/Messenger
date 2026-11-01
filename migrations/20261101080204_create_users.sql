-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(64) NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL
);

-- +goose Down
DROP TABLE users;
