-- name: CreateLike :exec
INSERT INTO likes (
    user_id,
    beat_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetLikeByUserAndBeat :one
SELECT * FROM likes
WHERE user_id = $1 AND beat_id = $2
LIMIT 1;

-- name: ListLikesByUser :many
SELECT * FROM likes
WHERE user_id = $1
LIMIT $2
OFFSET $3;

-- name: ListLikesByBeat :many
SELECT * FROM likes
WHERE beat_id = $1
LIMIT $2
OFFSET $3;

-- name: DeleteLike :exec
DELETE FROM likes
WHERE user_id = $1 AND beat_id = $2;