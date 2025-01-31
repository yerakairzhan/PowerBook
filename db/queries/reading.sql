-- name: CreateReadingLog :one
INSERT INTO reading_logs (userid, date, minutes_read)
VALUES ($1, $2, $3)
        RETURNING *;

-- name: UpdateReadingLog :one
update reading_logs set minutes_read = $3 where (userid = $1 and date = $2)
    returning *;

-- name: GetReadingLogsByUser :many
SELECT date, minutes_read
FROM reading_logs
WHERE userid = $1
ORDER BY date DESC;

-- name: GetTopReaders :many
SELECT u.username, SUM(r.minutes_read) AS total_minutes
FROM reading_logs r
         JOIN users u ON r.userid = u.id
GROUP BY u.username
ORDER BY total_minutes DESC
    LIMIT $1;
