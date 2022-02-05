package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// executes a function within a db transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// contains input parameters for LikeBeatTx function
type LikeBeatTxParams struct {
	UserId int32 `json:"user_id"`
	BeatId int32 `json:"beat_id"`
}

// contains result of LikeBeatTx function
type LikeBeatTxResult struct {
	LikeEntry Like `json:"like_entry"`
	LikeUser  User `json:"like_user"`
	LikeBeat  Beat `json:"like_beat"`
}
