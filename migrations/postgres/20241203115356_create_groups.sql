-- +goose Up
-- +goose StatementBegin

CREATE TABLE groups(
    id SERIAL PRIMARY KEY,
    title VARCHAR(30) NOT NULL,
    link VARCHAR(25) UNIQUE NOT NULL,
    owner_id BIGINT REFERENCES users(id)
);

CREATE TABLE group_profiles(
    g_id BIGINT REFERENCES groups(id),
    link TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY(g_id , created_at)
);

CREATE TABLE group_members(
    group_id INT REFERENCES groups(id),
    member_id INT REFERENCES users(id),
    joined_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP,
    PRIMARY KEY(group_id,member_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE group_members;

DROP TABLE group_profiles;

DROP TABLE groups;
-- +goose StatementEnd
