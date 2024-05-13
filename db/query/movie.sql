-- name: GetMovie :one
SELECT *
FROM movies
WHERE id = ?
LIMIT 1;
-- name: CreateMovie :execresult
INSERT INTO movies (
        id,
        title,
        backdrop_path,
        poster_path,
        vote_average
    )
VALUES (?, ?, ?, ?, ?);