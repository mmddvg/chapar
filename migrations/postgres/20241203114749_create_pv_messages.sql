-- +goose Up
-- +goose StatementBegin

CREATE TABLE pv_messages(
    id BIGINT,
    pv_id INT REFERENCES private_chats(id),
    sender_id BIGINT REFERENCES users(id),
    message VARCHAR(250) NOT NULL,
    seen_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(pv_id,id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE pv_messages;

-- +goose StatementEnd
