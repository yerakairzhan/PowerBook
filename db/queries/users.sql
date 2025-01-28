-- name: CreateUser :one
INSERT INTO users (userid, username)
VALUES ($1, $2)
    RETURNING id, userid, username, registered, created_at;

-- name: GetLanguage :one
select language from users where userid = $1;

-- name: GetUser :one
SELECT userid, username FROM users WHERE userid = $1;