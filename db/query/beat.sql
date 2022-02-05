-- name: CreateBeat :one
INSERT INTO beats (
    creator_id,
    title,
    genre,
    key,
    bpm,
    tags,
    s3_key
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetBeatById :one
SELECT * FROM beats
WHERE id = $1
LIMIT 1;

-- name: ListBeatsById :many
SELECT * FROM beats
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListBeatsByCreatorId :many
SELECT * FROM beats
WHERE creator_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListBeatsByGenre :many
SELECT * FROM beats
WHERE genre = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListBeatsByKey :many
SELECT * FROM beats
WHERE key = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListBeatsByBpmRange :many
SELECT * FROM beats
WHERE bpm BETWEEN $1 AND $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: ListBeatsByCreatorIdAndGenre :many
SELECT * FROM beats
WHERE creator_id = $1 AND genre = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: ListBeatsByCreatorIdAndKey :many
SELECT * FROM beats
WHERE creator_id = $1 AND key = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: ListBeatsByCreatorIdAndBpmRange :many
SELECT * FROM beats
WHERE creator_id = $1 AND bpm BETWEEN $2 AND $3
ORDER BY id
LIMIT $4
OFFSET $5;

-- name: UpdateBeat :one
UPDATE beats
SET title = $2,
    genre = $3,
    key = $4,
    bpm = $5,
    tags = $6,
    s3_key = $7
WHERE id = $1
RETURNING *;

-- name: DeleteBeat :exec
DELETE FROM beats
WHERE id = $1;