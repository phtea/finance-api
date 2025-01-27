-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    balance DECIMAL(10, 2) NOT NULL DEFAULT 0
);
-- +goose Down
DROP TABLE users;
