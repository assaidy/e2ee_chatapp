-- +goose Up
-- +goose StatementBegin
create table sessions (
    id uuid default gen_random_uuid(),
    credentials_id uuid not null,
    token varchar not null unique,
    csrf_token varchar not null unique,
    created_at timestamptz not null default now(),

    primary key (id),
    foreign key (credentials_id) references credentials (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table sessions;
-- +goose StatementEnd
