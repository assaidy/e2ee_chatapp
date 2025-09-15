-- +goose Up
-- +goose StatementBegin
create table sessions (
    id uuid default gen_random_uuid(),
    user_id uuid not null,
    session_token varchar not null unique,
    csrf_token varchar not null unique,
    created_at timestamptz not null default now(),
    user_agent varchar not null,
    ip_address varchar not null,
    last_active timestamptz not null default now(),

    primary key (id),
    foreign key (user_id) references users (id) on delete cascade
);
-- +goose StatementEnd

-- +goose StatementBegin
create index on sessions(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table sessions;
-- +goose StatementEnd
