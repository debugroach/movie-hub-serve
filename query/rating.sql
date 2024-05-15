-- name: GetRatings :many
SELECT *
FROM ratings
WHERE username = ?;
-- name: GetRating :one
SELECT *
FROM ratings
WHERE username = ?
  AND movie_id = ?;
-- name: ListRatings :many
SELECT *
FROM ratings;
-- name: CreateRating :execresult
INSERT INTO ratings (username, movie_id, rating)
VALUES (?, ?, ?);
-- name: UpdateRating :execresult
UPDATE ratings
SET rating = ?
WHERE username = ?
  AND movie_id = ?;