-- name: CheckUsername :one
select exists (select 1 from users where username = $1 for update);

-- name: CheckEmail :one
select exists (select 1 from users where email = $1 for update);

-- name: InsertUser :exec
insert into users (id, name, username, email, password_hash)
values ($1, $2, $3, $4, $5);

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: GetUserByID :one
select * from users where id = $1;

-- name: UpdateUser :exec
update users 
set
    name = $1,
    username = $2,
    email = $3,
    password_hash = $4
where id = $5;

-- name: DeleteUserByID :exec
delete from users where id = $1;

-- name: InsertEmailVerificationToken :exec
insert into email_verification_tokens (id, user_id, created_at, expires_at)
values ($1, $2, $3, $4);

-- name: DeleteExpiredEmailVerificationTokens :exec
delete from email_verification_tokens where expires_at <= now();

-- name: GetEmailVerificationTokenByID :one
select * from email_verification_tokens where id = $1;

-- name: MarkEmailAsVerified :exec
update users set email_is_verified = true where id = $1;
