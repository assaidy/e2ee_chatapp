-- +goose Up
-- +goose StatementBegin
create table credentials (
    id uuid default gen_random_uuid(),
    email varchar(255) not null unique,
    email_is_verified bool not null default false,
    password_hash varchar(100) not null,
    created_at timestamptz not null default now(),

    primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table credentials;
-- +goose StatementEnd
