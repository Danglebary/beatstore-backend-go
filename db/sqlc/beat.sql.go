// Code generated by sqlc. DO NOT EDIT.
// source: beat.sql

package db

import (
	"context"
)

const createBeat = `-- name: CreateBeat :one
INSERT INTO beats (
    creator_id,
    title,
    genre,
    key,
    bpm,
    tags,
    s3_key,
    likes_count
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at
`

type CreateBeatParams struct {
	CreatorID  int32  `json:"creator_id"`
	Title      string `json:"title"`
	Genre      string `json:"genre"`
	Key        string `json:"key"`
	Bpm        int16  `json:"bpm"`
	Tags       string `json:"tags"`
	S3Key      string `json:"s3_key"`
	LikesCount int64  `json:"likes_count"`
}

func (q *Queries) CreateBeat(ctx context.Context, arg CreateBeatParams) (Beat, error) {
	row := q.db.QueryRowContext(ctx, createBeat,
		arg.CreatorID,
		arg.Title,
		arg.Genre,
		arg.Key,
		arg.Bpm,
		arg.Tags,
		arg.S3Key,
		arg.LikesCount,
	)
	var i Beat
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Genre,
		&i.Key,
		&i.Bpm,
		&i.Tags,
		&i.S3Key,
		&i.LikesCount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteBeat = `-- name: DeleteBeat :exec
DELETE FROM beats
WHERE id = $1
`

func (q *Queries) DeleteBeat(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteBeat, id)
	return err
}

const getBeatById = `-- name: GetBeatById :one
SELECT id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at FROM beats
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetBeatById(ctx context.Context, id int32) (Beat, error) {
	row := q.db.QueryRowContext(ctx, getBeatById, id)
	var i Beat
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Genre,
		&i.Key,
		&i.Bpm,
		&i.Tags,
		&i.S3Key,
		&i.LikesCount,
		&i.CreatedAt,
	)
	return i, err
}

const listBeatsByBpmRange = `-- name: ListBeatsByBpmRange :many
SELECT id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at FROM beats
WHERE bpm BETWEEN $1 AND $2
ORDER BY id
LIMIT $3
OFFSET $4
`

type ListBeatsByBpmRangeParams struct {
	Bpm    int16 `json:"bpm"`
	Bpm_2  int16 `json:"bpm_2"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListBeatsByBpmRange(ctx context.Context, arg ListBeatsByBpmRangeParams) ([]Beat, error) {
	rows, err := q.db.QueryContext(ctx, listBeatsByBpmRange,
		arg.Bpm,
		arg.Bpm_2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Beat
	for rows.Next() {
		var i Beat
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.Genre,
			&i.Key,
			&i.Bpm,
			&i.Tags,
			&i.S3Key,
			&i.LikesCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBeatsByCreatorId = `-- name: ListBeatsByCreatorId :many
SELECT id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at FROM beats
WHERE creator_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListBeatsByCreatorIdParams struct {
	CreatorID int32 `json:"creator_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListBeatsByCreatorId(ctx context.Context, arg ListBeatsByCreatorIdParams) ([]Beat, error) {
	rows, err := q.db.QueryContext(ctx, listBeatsByCreatorId, arg.CreatorID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Beat
	for rows.Next() {
		var i Beat
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.Genre,
			&i.Key,
			&i.Bpm,
			&i.Tags,
			&i.S3Key,
			&i.LikesCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBeatsByCreatorIdAndBpmRange = `-- name: ListBeatsByCreatorIdAndBpmRange :many
SELECT id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at FROM beats
WHERE creator_id = $1 AND bpm BETWEEN $2 AND $3
ORDER BY id
LIMIT $4
OFFSET $5
`

type ListBeatsByCreatorIdAndBpmRangeParams struct {
	CreatorID int32 `json:"creator_id"`
	Bpm       int16 `json:"bpm"`
	Bpm_2     int16 `json:"bpm_2"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListBeatsByCreatorIdAndBpmRange(ctx context.Context, arg ListBeatsByCreatorIdAndBpmRangeParams) ([]Beat, error) {
	rows, err := q.db.QueryContext(ctx, listBeatsByCreatorIdAndBpmRange,
		arg.CreatorID,
		arg.Bpm,
		arg.Bpm_2,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Beat
	for rows.Next() {
		var i Beat
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.Genre,
			&i.Key,
			&i.Bpm,
			&i.Tags,
			&i.S3Key,
			&i.LikesCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBeatsByCreatorIdAndGenre = `-- name: ListBeatsByCreatorIdAndGenre :many
SELECT id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at FROM beats
WHERE creator_id = $1 AND genre = $2
ORDER BY id
LIMIT $3
OFFSET $4
`

type ListBeatsByCreatorIdAndGenreParams struct {
	CreatorID int32  `json:"creator_id"`
	Genre     string `json:"genre"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

func (q *Queries) ListBeatsByCreatorIdAndGenre(ctx context.Context, arg ListBeatsByCreatorIdAndGenreParams) ([]Beat, error) {
	rows, err := q.db.QueryContext(ctx, listBeatsByCreatorIdAndGenre,
		arg.CreatorID,
		arg.Genre,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Beat
	for rows.Next() {
		var i Beat
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.Genre,
			&i.Key,
			&i.Bpm,
			&i.Tags,
			&i.S3Key,
			&i.LikesCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBeatsByCreatorIdAndKey = `-- name: ListBeatsByCreatorIdAndKey :many
SELECT id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at FROM beats
WHERE creator_id = $1 AND key = $2
ORDER BY id
LIMIT $3
OFFSET $4
`

type ListBeatsByCreatorIdAndKeyParams struct {
	CreatorID int32  `json:"creator_id"`
	Key       string `json:"key"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

func (q *Queries) ListBeatsByCreatorIdAndKey(ctx context.Context, arg ListBeatsByCreatorIdAndKeyParams) ([]Beat, error) {
	rows, err := q.db.QueryContext(ctx, listBeatsByCreatorIdAndKey,
		arg.CreatorID,
		arg.Key,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Beat
	for rows.Next() {
		var i Beat
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.Genre,
			&i.Key,
			&i.Bpm,
			&i.Tags,
			&i.S3Key,
			&i.LikesCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBeatsByGenre = `-- name: ListBeatsByGenre :many
SELECT id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at FROM beats
WHERE genre = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListBeatsByGenreParams struct {
	Genre  string `json:"genre"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListBeatsByGenre(ctx context.Context, arg ListBeatsByGenreParams) ([]Beat, error) {
	rows, err := q.db.QueryContext(ctx, listBeatsByGenre, arg.Genre, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Beat
	for rows.Next() {
		var i Beat
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.Genre,
			&i.Key,
			&i.Bpm,
			&i.Tags,
			&i.S3Key,
			&i.LikesCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBeatsById = `-- name: ListBeatsById :many
SELECT id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at FROM beats
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListBeatsByIdParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListBeatsById(ctx context.Context, arg ListBeatsByIdParams) ([]Beat, error) {
	rows, err := q.db.QueryContext(ctx, listBeatsById, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Beat
	for rows.Next() {
		var i Beat
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.Genre,
			&i.Key,
			&i.Bpm,
			&i.Tags,
			&i.S3Key,
			&i.LikesCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBeatsByKey = `-- name: ListBeatsByKey :many
SELECT id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at FROM beats
WHERE key = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListBeatsByKeyParams struct {
	Key    string `json:"key"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListBeatsByKey(ctx context.Context, arg ListBeatsByKeyParams) ([]Beat, error) {
	rows, err := q.db.QueryContext(ctx, listBeatsByKey, arg.Key, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Beat
	for rows.Next() {
		var i Beat
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.Title,
			&i.Genre,
			&i.Key,
			&i.Bpm,
			&i.Tags,
			&i.S3Key,
			&i.LikesCount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBeat = `-- name: UpdateBeat :one
UPDATE beats
SET title = $2,
    genre = $3,
    key = $4,
    bpm = $5,
    tags = $6,
    s3_key = $7,
    likes_count = $8
WHERE id = $1
RETURNING id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at
`

type UpdateBeatParams struct {
	ID         int32  `json:"id"`
	Title      string `json:"title"`
	Genre      string `json:"genre"`
	Key        string `json:"key"`
	Bpm        int16  `json:"bpm"`
	Tags       string `json:"tags"`
	S3Key      string `json:"s3_key"`
	LikesCount int64  `json:"likes_count"`
}

func (q *Queries) UpdateBeat(ctx context.Context, arg UpdateBeatParams) (Beat, error) {
	row := q.db.QueryRowContext(ctx, updateBeat,
		arg.ID,
		arg.Title,
		arg.Genre,
		arg.Key,
		arg.Bpm,
		arg.Tags,
		arg.S3Key,
		arg.LikesCount,
	)
	var i Beat
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Genre,
		&i.Key,
		&i.Bpm,
		&i.Tags,
		&i.S3Key,
		&i.LikesCount,
		&i.CreatedAt,
	)
	return i, err
}

const updateBeatLikesCount = `-- name: UpdateBeatLikesCount :one
UPDATE beats
SET likes_count = $2
WHERE id = $1
RETURNING id, creator_id, title, genre, key, bpm, tags, s3_key, likes_count, created_at
`

type UpdateBeatLikesCountParams struct {
	ID         int32 `json:"id"`
	LikesCount int64 `json:"likes_count"`
}

func (q *Queries) UpdateBeatLikesCount(ctx context.Context, arg UpdateBeatLikesCountParams) (Beat, error) {
	row := q.db.QueryRowContext(ctx, updateBeatLikesCount, arg.ID, arg.LikesCount)
	var i Beat
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Title,
		&i.Genre,
		&i.Key,
		&i.Bpm,
		&i.Tags,
		&i.S3Key,
		&i.LikesCount,
		&i.CreatedAt,
	)
	return i, err
}
