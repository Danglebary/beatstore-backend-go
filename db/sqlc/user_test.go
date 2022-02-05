package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/danglebary/beatstore-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
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

func deleteRandomUser(t *testing.T, id int32) {
	err := testQueries.DeleteUser(context.Background(), id)
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	user := createRandomUser(t)
	deleteRandomUser(t, user.ID)
}

func TestGetUserById(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserById(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	deleteRandomUser(t, user1.ID)
}

func TestGetUserByUsername(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUserByUsername(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

	deleteRandomUser(t, user1.ID)
}

func TestListUsersById(t *testing.T) {
	n := 10
	for i := 0; i < n; i++ {
		createRandomUser(t)
	}
	arg := ListUsersByIdParams{
		Limit:  5,
		Offset: 5,
	}
	users, err := testQueries.ListUsersById(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 5, len(users))

	for _, user := range users {
		require.NotEmpty(t, user)
		deleteRandomUser(t, user.ID)
	}
}

func TestListUsersByUsername(t *testing.T) {
	n := 10
	for i := 0; i < n; i++ {
		createRandomUser(t)
	}
	arg := ListUsersByUsernameParams{
		Limit:  5,
		Offset: 5,
	}
	users, err := testQueries.ListUsersByUsername(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 5, len(users))

	for _, user := range users {
		require.NotEmpty(t, user)
		deleteRandomUser(t, user.ID)
	}
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)
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

	deleteRandomUser(t, user1.ID)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserById(context.Background(), user1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user2)
}
