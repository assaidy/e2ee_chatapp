-- +goose Up
-- +goose StatementBegin
create table users (
    id uuid default gen_random_uuid(),
    name varchar(100) not null,
    username varchar(50) not null unique,
    email varchar(254) not null unique,
    password_hash varchar(100) not null,
    email_is_verified bool not null default false,
    joined_at timestamptz not null default now(),

    primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
