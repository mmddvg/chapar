-- +goose Up
-- +goose StatementBegin

CREATE TABLE private_chats(
    id SERIAL PRIMARY KEY,
    user1 BIGINT REFERENCES users(id),
    user2 BIGINT REFERENCES users(id),
    is_blocked BOOLEAN DEFAULT 'f',
    created_at TIMESTAMP DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE private_chats;

-- +goose StatementEnd
