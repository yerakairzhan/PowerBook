-- name: CreateUser :one
INSERT INTO users (userid, username)
VALUES ($1, $2)
    RETURNING id, userid, username, registered, created_at;

-- name: GetLanguage :one
select language from users where userid = $1;

-- name: GetUser :one
SELECT userid, username FROM users WHERE userid = $1;

-- name: SetUserState :exec
INSERT INTO users (userid, state)
VALUES ($1, $2)
    ON CONFLICT (userid) DO UPDATE
    SET state = EXCLUDED.state;

-- name: GetUserState :one
SELECT state FROM users WHERE userid = $1;

-- name: DeleteUserState :exec
update users set state = null where userid = $1;

-- name: SetRegistered :exec
update users set registered = true where userid = $1;

-- name: GetRegistered :one
select registered from users where userid = $1 ;