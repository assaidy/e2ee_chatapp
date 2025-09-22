-- name: CheckUsername :one
select exists (select 1 from users where username = $1 for update);

-- name: InsertUser :exec
insert into users (id, name, username, credentials_id)
values ($1, $2, $3, $4);
