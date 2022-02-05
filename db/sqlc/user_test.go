package db

import (
	"context"
	"testing"

	"github.com/danglebary/beatstore-backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		Username: util.RandomUserName(),
		Password: util.RandomPassword(),
		Email:    util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
}

func RandomUserName() {
	panic("unimplemented")
}
