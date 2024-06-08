-- +goose Up
create table users (
    id UUID PRIMARY KEY,
    create_at  TIMESTAMP NOT NULL,
    update_at  TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    user_name varchar(50) NOT NULL unique,
    password TEXT NOT NULL,
    role varchar(20) NOT NULL,
    api_key UUID,
    parent_user_id UUID REFERENCES users(id) ON DELETE CASCADE

);
-- +goose Down

DROP TABLE users;