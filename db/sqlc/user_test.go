package db

import (
	"context"
	"testing"
	"time"

	"github.com/alrasyidin/simplebank-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	getUser, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, getUser)

	require.Equal(t, user.Username, getUser.Username)
	require.Equal(t, user.Email, getUser.Email)
	require.Equal(t, user.HashedPassword, getUser.HashedPassword)
	require.Equal(t, user.FullName, getUser.FullName)
	require.WithinDuration(t, user.CreatedAt, getUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, getUser.PasswordChangedAt, time.Second)
}
