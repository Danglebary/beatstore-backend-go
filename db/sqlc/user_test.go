package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/danglebary/beatstore-backend/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandomUsername(),
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

	return user
}

func DeleteRandomUser(t *testing.T, id int32) {
	err := testQueries.DeleteUser(context.Background(), id)
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	user := CreateRandomUser(t)
	DeleteRandomUser(t, user.ID)
}

func TestGetUserById(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUserById(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	DeleteRandomUser(t, user1.ID)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	DeleteRandomUser(t, user1.ID)
}

func TestListUsersById(t *testing.T) {
	n := 10
	for i := 0; i < n; i++ {
		CreateRandomUser(t)
	}
	arg := ListUsersByIdParams{
		Limit:  int32(n),
		Offset: 0,
	}
	users, err := testQueries.ListUsersById(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, n, len(users))

	for _, user := range users {
		require.NotEmpty(t, user)
		DeleteRandomUser(t, user.ID)
	}
}

func TestListUsersByUsername(t *testing.T) {
	n := 10
	for i := 0; i < n; i++ {
		CreateRandomUser(t)
	}
	arg := ListUsersByUsernameParams{
		Limit:  int32(n),
		Offset: 0,
	}
	users, err := testQueries.ListUsersByUsername(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, n, len(users))

	for _, user := range users {
		require.NotEmpty(t, user)
		DeleteRandomUser(t, user.ID)
	}
}

func TestUpdateUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	arg := UpdateUserParams{
		ID:       user1.ID,
		Username: util.RandomUsername(),
		Password: util.RandomPassword(),
		Email:    util.RandomEmail(),
	}
	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, arg.Username, user2.Username)
	require.Equal(t, arg.Password, user2.Password)
	require.Equal(t, arg.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	DeleteRandomUser(t, user1.ID)
}

func TestDeleteUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserById(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}
