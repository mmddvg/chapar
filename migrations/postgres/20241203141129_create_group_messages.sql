-- +goose Up
-- +goose StatementBegin

CREATE TABLE group_messages(
    id BIGINT PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES groups(id),
    message VARCHAR(250) NOT NULL,
    sender_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT now()
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE group_messages;

-- +goose StatementEnd
