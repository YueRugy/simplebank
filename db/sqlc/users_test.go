package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"
)

func TestQueries_CreateUser(t *testing.T) {
	pwd, err := util.HashPassword(util.RandomOwner())
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: pwd,
		Email:          util.RandomOwner(),
		FullName:       util.RandomOwner(),
	}
	user, err := testQuery.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.NotZero(t, user.CreateAt)
}
func createRandomUser(t *testing.T) User {
	pwd, err := util.HashPassword(util.RandomOwner())
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: pwd,
		Email:          util.RandomOwner(),
		FullName:       util.RandomOwner(),
	}
	user, err := testQuery.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.NotZero(t, user.CreateAt)
	require.NotZero(t, user.PasswordChangedAt)
	return user
}

func TestQueries_GetUsers(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQuery.GetUsers(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreateAt, user2.CreateAt, 1*time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, 1*time.Second)
}
