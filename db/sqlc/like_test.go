package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateLike(t *testing.T) {
	user1 := createRandomUser(t)
	beat1 := createRandomBeat(t)

	like, err := testQueries.CreateLike(context.Background(), CreateLikeParams{
		UserID: user1.ID,
		BeatID: beat1.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, like)
	require.Equal(t, user1.ID, like.UserID)
	require.Equal(t, beat1.ID, like.BeatID)

	err = testQueries.DeleteLike(context.Background(), DeleteLikeParams{
		UserID: user1.ID,
		BeatID: beat1.ID,
	})
	require.NoError(t, err)

	deleteRandomBeat(t, beat1.ID)
	deleteRandomUser(t, user1.ID)
	deleteRandomUser(t, beat1.CreatorID)
}

func TestListLikesByUser(t *testing.T) {
	user1 := createRandomUser(t)

	n := 10
	for i := 0; i < n; i++ {
		beat := createRandomBeat(t)
		like, err := testQueries.CreateLike(context.Background(), CreateLikeParams{
			UserID: user1.ID,
			BeatID: beat.ID,
		})
		require.NoError(t, err)
		require.NotEmpty(t, like)
		require.Equal(t, user1.ID, like.UserID)
		require.Equal(t, beat.ID, like.BeatID)
	}

	likes, err := testQueries.ListLikesByUser(context.Background(), ListLikesByUserParams{
		UserID: user1.ID,
		Limit:  int32(n),
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, likes)

	for _, like := range likes {
		require.NotEmpty(t, like)
		require.Equal(t, user1.ID, like.UserID)

		err = testQueries.DeleteLike(context.Background(), DeleteLikeParams{
			UserID: user1.ID,
			BeatID: like.BeatID,
		})
		require.NoError(t, err)
		beat, err := testQueries.GetBeatById(context.Background(), like.BeatID)
		require.NoError(t, err)
		require.NotEmpty(t, beat)

		deleteRandomBeat(t, like.BeatID)
		deleteRandomUser(t, beat.CreatorID)
	}

	deleteRandomUser(t, user1.ID)
}

func TestListLikesByBeat(t *testing.T) {
	beat1 := createRandomBeat(t)

	n := 10
	for i := 0; i < n; i++ {
		user := createRandomUser(t)
		like, err := testQueries.CreateLike(context.Background(), CreateLikeParams{
			UserID: user.ID,
			BeatID: beat1.ID,
		})
		require.NoError(t, err)
		require.NotEmpty(t, like)
		require.Equal(t, user.ID, like.UserID)
		require.Equal(t, beat1.ID, like.BeatID)
	}

	likes, err := testQueries.ListLikesByBeat(context.Background(), ListLikesByBeatParams{
		BeatID: beat1.ID,
		Limit:  int32(n),
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, likes)

	for _, like := range likes {
		require.NotEmpty(t, like)
		require.Equal(t, beat1.ID, like.BeatID)

		err = testQueries.DeleteLike(context.Background(), DeleteLikeParams{
			UserID: like.UserID,
			BeatID: beat1.ID,
		})
		require.NoError(t, err)
		user, err := testQueries.GetUserById(context.Background(), like.UserID)
		require.NoError(t, err)
		require.NotEmpty(t, user)

		deleteRandomUser(t, user.ID)
	}
	deleteRandomBeat(t, beat1.ID)
	deleteRandomUser(t, beat1.CreatorID)
}

func TestDeleteLike(t *testing.T) {
	user1 := createRandomUser(t)
	beat1 := createRandomBeat(t)

	like, err := testQueries.CreateLike(context.Background(), CreateLikeParams{
		UserID: user1.ID,
		BeatID: beat1.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, like)
	require.Equal(t, user1.ID, like.UserID)
	require.Equal(t, beat1.ID, like.BeatID)

	err = testQueries.DeleteLike(context.Background(), DeleteLikeParams{
		UserID: user1.ID,
		BeatID: beat1.ID,
	})
	require.NoError(t, err)

	like, err = testQueries.GetLikeByUserAndBeat(context.Background(), GetLikeByUserAndBeatParams{
		UserID: user1.ID,
		BeatID: beat1.ID,
	})
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, like)

	deleteRandomBeat(t, beat1.ID)
	deleteRandomUser(t, beat1.CreatorID)
	deleteRandomUser(t, user1.ID)
}
