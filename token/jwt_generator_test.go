package token

import (
	"testing"
	"time"

	"github.com/alrasyidin/simplebank-go/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWTGenerator(t *testing.T) {
	generator, err := NewJWTGenerator(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := EXPIRE_TOKEN

	issuedAt := time.Now()

	token, err := generator.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := generator.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.Payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt.Time, time.Second)
	require.WithinDuration(t, issuedAt.Add(EXPIRE_TOKEN), payload.ExpiresAt.Time, time.Second)

}

func TestJWTGeneratorExpired(t *testing.T) {
	generator, err := NewJWTGenerator(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := -time.Minute

	token, err := generator.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := generator.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayloadJWT(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	generator, err := NewJWTGenerator(util.RandomString(32))
	require.NoError(t, err)
	payload, err = generator.VerifyToken(token)

	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
