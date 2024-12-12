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
    PRIMARY KEY(group_id,member_id)
);

CREATE OR REPLACE FUNCTION add_owner_to_group_members()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO group_members (group_id, member_id, joined_at)
    VALUES (NEW.id, NEW.owner_id, NOW());
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_group_insert
AFTER INSERT ON groups
FOR EACH ROW
EXECUTE FUNCTION add_owner_to_group_members();


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS after_group_insert ON groups;

DROP FUNCTION IF EXISTS add_owner_to_group_members();

DROP TABLE group_members;

DROP TABLE group_profiles;

DROP TABLE groups;
-- +goose StatementEnd
