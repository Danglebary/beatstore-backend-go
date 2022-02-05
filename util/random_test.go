package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	var min int64 = 1
	var max int64 = 10

	res := RandomInt(min, max)
	require.GreaterOrEqual(t, res, min)
	require.LessOrEqual(t, res, max)
}

func TestRandomString(t *testing.T) {
	var k int = 6
	res := RandomString(k)
	require.Equal(t, k, len(res))
}

func TestRandomUsername(t *testing.T) {
	res := RandomUsername()
	require.Equal(t, 6, len(res))
}

func TestRandomPassword(t *testing.T) {
	res := RandomPassword()
	require.Equal(t, 6, len(res))
}

func TestRandomEmail(t *testing.T) {
	res := RandomEmail()
	require.Contains(t, res, "@")
	require.Contains(t, res, ".com")
}

func TestRandomTitle(t *testing.T) {
	res := RandomTitle()
	require.Equal(t, 6, len(res))
}

func TestRandomGenre(t *testing.T) {
	res := RandomGenre()
	require.Equal(t, 6, len(res))
}

func TestRandomKey(t *testing.T) {
	res := RandomKey()
	require.Contains(t, MusicalKeys, res)
}

func TestRandomBpm(t *testing.T) {
	res := RandomBpm()
	require.GreaterOrEqual(t, res, int16(20))
	require.LessOrEqual(t, res, int16(999))
}

func TestRandomTags(t *testing.T) {
	res := RandomTags()
	require.Contains(t, res, ",")
}

func TestRandomS3Key(t *testing.T) {
	res := RandomS3Key()
	require.Equal(t, 12, len(res))
}

func TestRandomLikesCount(t *testing.T) {
	res := RandomLikesCount()
	require.GreaterOrEqual(t, res, int64(0))
	require.LessOrEqual(t, res, int64(1000))
}
