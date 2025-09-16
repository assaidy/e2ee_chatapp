-- +goose Up
-- +goose StatementBegin
create table email_verification_tokens (
    id uuid,
    user_id uuid not null,
    created_at timestamptz not null default now(),
    expires_at timestamptz not null,

    primary key (id),
    foreign key (user_id) references users (id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table email_verification_tokens;
-- +goose StatementEnd
