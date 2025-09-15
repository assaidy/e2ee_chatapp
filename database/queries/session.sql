-- name: InsertSession :one
insert into sessions (id, user_id, session_token, csrf_token, user_agent, ip_address)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: GetSessionByID :one
select * from sessions where id = $1;

-- name: DeleteSessionForUser :execrows
delete from sessions where id = $1 and user_id = $2;

-- name: UpdateSessionLastActive :exec
update sessions set last_active = now() where id = $1;
