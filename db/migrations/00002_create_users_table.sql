-- +goose Up
-- +goose StatementBegin
create table users (
    id uuid default gen_random_uuid(),
    name varchar(50) not null,
    username varchar(50) not null unique,
    credentials_id uuid not null unique,
    created_at timestamptz not null default now(),

    primary key (id),
    foreign key (credentials_id) references credentials (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
