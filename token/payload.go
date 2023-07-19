package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrExpiredToken = errors.New("token is expired")
var ErrInvalidToken = errors.New("token is invalid")

const EXPIRE_TOKEN = time.Hour * 24 * 7

type PayloadClaim struct {
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type PayloadJWT struct {
	Payload
	jwt.RegisteredClaims
}

func NewPayloadJWT(username string, expire time.Duration) (*PayloadJWT, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &PayloadJWT{
		Payload: Payload{
			ID:       tokenId,
			Username: username,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}, nil
}

type PayloadPaseto struct {
	Payload
	IssuedAt  time.Time
	ExpiresAt time.Time
}

func NewPayloadPaseto(username string, expire time.Duration) (*PayloadPaseto, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &PayloadPaseto{
		Payload: Payload{
			ID:       tokenId,
			Username: username,
		},
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(expire),
	}, nil
}

func (payload *PayloadPaseto) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}

	return nil
}
