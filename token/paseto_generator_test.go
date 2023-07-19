package token

import (
	"testing"
	"time"

	"github.com/alrasyidin/simplebank-go/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoGenerator(t *testing.T) {
	// payload, err := NewPayloadPaseto(util.RandomOwner(), time.Minute)
	// require.NoError(t, err)

	generator, err := NewPasetoGenerator(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := generator.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := generator.VerifyTokenPaseto(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

func TestPasetoGeneratorExpired(t *testing.T) {
	generator, err := NewPasetoGenerator(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, err := generator.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := generator.VerifyTokenPaseto(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
