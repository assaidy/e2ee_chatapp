-- name: CheckEmail :one
select exists (select 1 from credentials where email = $1 for update);

-- name: InsertCredentials :exec
insert into credentials (id, email, password_hash)
values ($1, $2, $3);

-- name: InsertEmailVerificationToken :exec
insert into email_verification_tokens (id, email, expires_at)
values ($1, $2, $3);

-- name: DeleteStaleEmailVerificationTokens :exec
delete from email_verification_tokens where expires_at <= now();

-- name: GetEmailVerificationTokenByID :one
select * from email_verification_tokens where id = $1;

-- name: MarkEmailAsVerified :exec
update credentials set email_is_verified = true where email = $1;

-- name: GetCredentialsByEmail :one
select * from credentials where email = $1;

-- name: InsertSession :one
insert into sessions (id, credentials_id, token, csrf_token)
values ($1, $2, $3, $4)
returning *;

-- name: GetSessionByID :one
select * from sessions where id = $1;
