-- +goose Up
-- +goose StatementBegin

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    username VARCHAR(25) UNIQUE NOT NULL,
    name VARCHAR(30) NOT NULL
);

CREATE TABLE user_profiles(
    user_id BIGINT REFERENCES users(id),
    link TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(user_id,created_at) 
);


CREATE TABLE contacts(
    user_id BIGINT REFERENCES users(id),
    contact_id BIGINT REFERENCES users(id),
    PRIMARY KEY(user_id,contact_id)
);

CREATE TABLE blocked(
    user_id BIGINT REFERENCES users(id),
    target_id BIGINT REFERENCES users(id),
    CONSTRAINT cant_block_self CHECK (user_id != target_id) 
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE blocked;

DROP TABLE contacts;

DROP TABLE user_profiles;

DROP TABLE users;

-- +goose StatementEnd
