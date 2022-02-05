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
