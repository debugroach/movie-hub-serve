-- name: GetUser :one
SELECT *
FROM users
WHERE username = ?
LIMIT 1;
-- name: ListUsers :many
SELECT *
FROM users
LIMIT ? OFFSET ?;
-- name: CreateUser :execresult
INSERT INTO users (username, password)
VALUES (?, ?);
-- name: DeleteUser :exec
DELETE FROM users
WHERE username = ?;