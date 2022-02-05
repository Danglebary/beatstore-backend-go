package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/danglebary/beatstore-backend-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomBeat(t *testing.T) Beat {
	user1 := createRandomUser(t)

	arg := CreateBeatParams{
		CreatorID: user1.ID,
		Title:     util.RandomTitle(),
		Genre:     util.RandomGenre(),
		Key:       util.RandomKey(),
		Bpm:       util.RandomBpm(),
		Tags:      util.RandomTags(),
		S3Key:     util.RandomS3Key(),
	}

	beat, err := testQueries.CreateBeat(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beat)

	require.Equal(t, user1.ID, beat.CreatorID)
	require.Equal(t, arg.Title, beat.Title)
	require.Equal(t, arg.Genre, beat.Genre)
	require.Equal(t, arg.Key, beat.Key)
	require.Equal(t, arg.Bpm, beat.Bpm)
	require.Equal(t, arg.Tags, beat.Tags)
	require.Equal(t, arg.S3Key, beat.S3Key)

	require.NotZero(t, beat.ID)
	require.NotZero(t, beat.CreatedAt)

	return beat
}

func createRandomBeatWithArgs(t *testing.T, arg CreateBeatParams) Beat {
	beat, err := testQueries.CreateBeat(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beat)

	require.Equal(t, arg.CreatorID, beat.CreatorID)
	require.Equal(t, arg.Title, beat.Title)
	require.Equal(t, arg.Genre, beat.Genre)
	require.Equal(t, arg.Key, beat.Key)
	require.Equal(t, arg.Bpm, beat.Bpm)
	require.Equal(t, arg.Tags, beat.Tags)
	require.Equal(t, arg.S3Key, beat.S3Key)

	require.NotZero(t, beat.ID)
	require.NotZero(t, beat.CreatedAt)

	return beat
}

func deleteRandomBeat(t *testing.T, id int32) {
	err := testQueries.DeleteBeat(context.Background(), id)
	require.NoError(t, err)
}

func TestCreateBeat(t *testing.T) {
	beat1 := createRandomBeat(t)

	deleteRandomBeat(t, beat1.ID)
	deleteRandomUser(t, beat1.CreatorID)
}

func TestUpdateBeat(t *testing.T) {
	beat1 := createRandomBeat(t)

	arg := UpdateBeatParams{
		ID:    beat1.ID,
		Title: util.RandomTitle(),
		Genre: util.RandomGenre(),
		Key:   util.RandomKey(),
		Bpm:   util.RandomBpm(),
		Tags:  util.RandomTags(),
		S3Key: util.RandomS3Key(),
	}

	beat2, err := testQueries.UpdateBeat(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beat2)

	require.Equal(t, beat1.ID, beat2.ID)
	require.Equal(t, arg.Title, beat2.Title)
	require.Equal(t, arg.Genre, beat2.Genre)
	require.Equal(t, arg.Key, beat2.Key)
	require.Equal(t, arg.Bpm, beat2.Bpm)
	require.Equal(t, arg.Tags, beat2.Tags)
	require.Equal(t, arg.S3Key, beat2.S3Key)
	require.WithinDuration(t, beat1.CreatedAt, beat2.CreatedAt, time.Second)

	deleteRandomBeat(t, beat1.ID)
	deleteRandomUser(t, beat1.CreatorID)
}

func TestGetBeatById(t *testing.T) {
	beat1 := createRandomBeat(t)

	beat2, err := testQueries.GetBeatById(context.Background(), beat1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, beat2)

	require.Equal(t, beat1.ID, beat2.ID)
	require.Equal(t, beat1.CreatorID, beat2.CreatorID)
	require.Equal(t, beat1.Title, beat2.Title)
	require.Equal(t, beat1.Genre, beat2.Genre)
	require.Equal(t, beat1.Key, beat2.Key)
	require.Equal(t, beat1.Bpm, beat2.Bpm)
	require.Equal(t, beat1.S3Key, beat2.S3Key)
	require.WithinDuration(t, beat1.CreatedAt, beat2.CreatedAt, time.Second)

	deleteRandomBeat(t, beat1.ID)
	deleteRandomUser(t, beat1.CreatorID)
}

func TestListBeatsById(t *testing.T) {
	n := 10
	for i := 0; i < n; i++ {
		createRandomBeat(t)
	}
	arg := ListBeatsByIdParams{
		Limit:  int32(n),
		Offset: 0,
	}
	beats, err := testQueries.ListBeatsById(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beats)
	require.Equal(t, n, len(beats))

	for _, beat := range beats {
		require.NotEmpty(t, beat)
		deleteRandomBeat(t, beat.ID)
		deleteRandomUser(t, beat.CreatorID)
	}
}

func TestListBeatsByCreatorId(t *testing.T) {
	user1 := createRandomUser(t)

	beatArg := CreateBeatParams{
		CreatorID: user1.ID,
		Title:     util.RandomTitle(),
		Genre:     util.RandomGenre(),
		Key:       util.RandomKey(),
		Bpm:       util.RandomBpm(),
		Tags:      util.RandomTags(),
		S3Key:     util.RandomS3Key(),
	}

	n := 10
	for i := 0; i < n; i++ {
		createRandomBeatWithArgs(t, beatArg)
	}

	arg := ListBeatsByCreatorIdParams{
		CreatorID: user1.ID,
		Limit:     int32(n),
		Offset:    0,
	}
	beats, err := testQueries.ListBeatsByCreatorId(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beats)
	require.Equal(t, n, len(beats))

	for _, beat := range beats {
		require.NotEmpty(t, beat)
		deleteRandomBeat(t, beat.ID)
	}

	deleteRandomUser(t, user1.ID)
}

func TestListBeatsbyGenre(t *testing.T) {
	genre := util.RandomGenre()

	n := 10
	for i := 0; i < n; i++ {
		user := createRandomUser(t)
		args := CreateBeatParams{
			CreatorID: user.ID,
			Title:     util.RandomTitle(),
			Genre:     genre,
			Key:       util.RandomKey(),
			Bpm:       util.RandomBpm(),
			Tags:      util.RandomTags(),
			S3Key:     util.RandomS3Key(),
		}
		createRandomBeatWithArgs(t, args)
	}

	arg := ListBeatsByGenreParams{
		Genre:  genre,
		Limit:  int32(n),
		Offset: 0,
	}
	beats, err := testQueries.ListBeatsByGenre(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beats)
	require.Equal(t, n, len(beats))

	for _, beat := range beats {
		require.NotEmpty(t, beat)
		require.Equal(t, genre, beat.Genre)
		deleteRandomBeat(t, beat.ID)
		deleteRandomUser(t, beat.CreatorID)
	}
}

func TestListBeatsByKey(t *testing.T) {
	key := util.RandomKey()

	n := 10
	for i := 0; i < n; i++ {
		user := createRandomUser(t)
		args := CreateBeatParams{
			CreatorID: user.ID,
			Title:     util.RandomTitle(),
			Genre:     util.RandomGenre(),
			Key:       key,
			Bpm:       util.RandomBpm(),
			Tags:      util.RandomTags(),
			S3Key:     util.RandomS3Key(),
		}
		createRandomBeatWithArgs(t, args)
	}

	arg := ListBeatsByKeyParams{
		Key:    key,
		Limit:  int32(n),
		Offset: 0,
	}
	beats, err := testQueries.ListBeatsByKey(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beats)
	require.Equal(t, n, len(beats))

	for _, beat := range beats {
		require.NotEmpty(t, beat)
		require.Equal(t, key, beat.Key)
		deleteRandomBeat(t, beat.ID)
		deleteRandomUser(t, beat.CreatorID)
	}
}

func TestListBeatsByBpmRange(t *testing.T) {
	min := 100
	max := 200

	n := 10
	for i := 0; i < n; i++ {
		user := createRandomUser(t)
		args := CreateBeatParams{
			CreatorID: user.ID,
			Title:     util.RandomTitle(),
			Genre:     util.RandomGenre(),
			Key:       util.RandomKey(),
			Bpm:       int16(util.RandomInt(int64(min), int64(max))),
			Tags:      util.RandomTags(),
			S3Key:     util.RandomS3Key(),
		}
		createRandomBeatWithArgs(t, args)
	}

	arg := ListBeatsByBpmRangeParams{
		Bpm:    int16(min),
		Bpm_2:  int16(max),
		Limit:  int32(n),
		Offset: 0,
	}
	beats, err := testQueries.ListBeatsByBpmRange(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beats)
	require.Equal(t, n, len(beats))

	for _, beat := range beats {
		require.NotEmpty(t, beat)
		require.GreaterOrEqual(t, beat.Bpm, int16(min))
		require.LessOrEqual(t, beat.Bpm, int16(max))
		deleteRandomBeat(t, beat.ID)
		deleteRandomUser(t, beat.CreatorID)
	}
}

func TestListBeatsByCreatorIdAndGenre(t *testing.T) {
	user1 := createRandomUser(t)
	genre := util.RandomGenre()

	n := 10
	for i := 0; i < n; i++ {
		args := CreateBeatParams{
			CreatorID: user1.ID,
			Title:     util.RandomTitle(),
			Genre:     genre,
			Key:       util.RandomKey(),
			Bpm:       util.RandomBpm(),
			Tags:      util.RandomTags(),
			S3Key:     util.RandomS3Key(),
		}
		createRandomBeatWithArgs(t, args)
	}

	arg := ListBeatsByCreatorIdAndGenreParams{
		CreatorID: user1.ID,
		Genre:     genre,
		Limit:     int32(n),
		Offset:    0,
	}
	beats, err := testQueries.ListBeatsByCreatorIdAndGenre(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beats)
	require.Equal(t, n, len(beats))

	for _, beat := range beats {
		require.NotEmpty(t, beat)
		require.Equal(t, user1.ID, beat.CreatorID)
		require.Equal(t, genre, beat.Genre)
		deleteRandomBeat(t, beat.ID)
	}
	deleteRandomUser(t, user1.ID)
}

func TestListBeatsByCreatorIdAndKey(t *testing.T) {
	user1 := createRandomUser(t)
	key := util.RandomKey()

	n := 10
	for i := 0; i < n; i++ {
		args := CreateBeatParams{
			CreatorID: user1.ID,
			Title:     util.RandomTitle(),
			Genre:     util.RandomGenre(),
			Key:       key,
			Bpm:       util.RandomBpm(),
			Tags:      util.RandomTags(),
			S3Key:     util.RandomS3Key(),
		}
		createRandomBeatWithArgs(t, args)
	}

	arg := ListBeatsByCreatorIdAndKeyParams{
		CreatorID: user1.ID,
		Key:       key,
		Limit:     int32(n),
		Offset:    0,
	}
	beats, err := testQueries.ListBeatsByCreatorIdAndKey(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beats)
	require.Equal(t, n, len(beats))

	for _, beat := range beats {
		require.NotEmpty(t, beat)
		require.Equal(t, user1.ID, beat.CreatorID)
		require.Equal(t, key, beat.Key)
		deleteRandomBeat(t, beat.ID)
	}
	deleteRandomUser(t, user1.ID)
}

func TestListBeatsByCreatorIdAndBpmRange(t *testing.T) {
	user1 := createRandomUser(t)
	min := 100
	max := 200

	n := 10
	for i := 0; i < n; i++ {
		args := CreateBeatParams{
			CreatorID: user1.ID,
			Title:     util.RandomTitle(),
			Genre:     util.RandomGenre(),
			Key:       util.RandomKey(),
			Bpm:       int16(util.RandomInt(int64(min), int64(max))),
			Tags:      util.RandomTags(),
			S3Key:     util.RandomS3Key(),
		}
		createRandomBeatWithArgs(t, args)
	}

	arg := ListBeatsByCreatorIdAndBpmRangeParams{
		CreatorID: user1.ID,
		Bpm:       int16(min),
		Bpm_2:     int16(max),
		Limit:     int32(n),
		Offset:    0,
	}
	beats, err := testQueries.ListBeatsByCreatorIdAndBpmRange(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, beats)
	require.Equal(t, n, len(beats))

	for _, beat := range beats {
		require.NotEmpty(t, beat)
		require.Equal(t, user1.ID, beat.CreatorID)
		require.GreaterOrEqual(t, beat.Bpm, int16(min))
		require.LessOrEqual(t, beat.Bpm, int16(max))
		deleteRandomBeat(t, beat.ID)
	}
	deleteRandomUser(t, user1.ID)
}

func TestDeleteBeat(t *testing.T) {
	beat1 := createRandomBeat(t)
	err := testQueries.DeleteBeat(context.Background(), beat1.ID)
	require.NoError(t, err)

	beat2, err := testQueries.GetBeatById(context.Background(), beat1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, beat2)

	deleteRandomUser(t, beat1.CreatorID)
}
