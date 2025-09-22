-- +goose Up
-- +goose StatementBegin
create table email_verification_tokens (
    id uuid,
    email varchar not null,
    created_at timestamptz not null default now(),
    expires_at timestamptz not null,

    primary key (id),
    foreign key (email) references credentials (email) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table email_verification_tokens;
-- +goose StatementEnd
