package db

import (
	"context"
	"database/sql"
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
func TestUpdateUserFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	fullName := util.RandomString(10)

	arg := UpdateUserUsingCaseSecondParams{
		FullName: sql.NullString{
			String: fullName,
			Valid:  true,
		},
		Username: oldUser.Username,
	}

	updatedUser, err := testQueries.UpdateUserUsingCaseSecond(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, fullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	email := util.RandomEmail()

	arg := UpdateUserUsingCaseSecondParams{
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Username: oldUser.Username,
	}

	updatedUser, err := testQueries.UpdateUserUsingCaseSecond(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, email, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateHashedPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := UpdateUserUsingCaseSecondParams{
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		Username: oldUser.Username,
	}

	updatedUser, err := testQueries.UpdateUserUsingCaseSecond(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, hashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullname := util.RandomOwner()
	newEmail := util.RandomEmail()

	newPassword := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := UpdateUserUsingCaseSecondParams{
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
		Username: oldUser.Username,
	}

	updatedUser, err := testQueries.UpdateUserUsingCaseSecond(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, hashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newFullname, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)
}
