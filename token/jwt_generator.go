package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTGenerator struct {
	secretKey string
}

// Constructor for JWTGenerator
func NewJWTGenerator(secretKey string) (Generator, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTGenerator{secretKey}, nil
}

func (g *JWTGenerator) CreateToken(username string, expire time.Duration) (string, error) {
	payload, err := NewPayloadJWT(username, expire)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(g.secretKey))
}

func (g *JWTGenerator) VerifyToken(token string) (*PayloadJWT, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &PayloadJWT{}, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(g.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*PayloadJWT)

	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func (mp *JWTGenerator) VerifyTokenPaseto(token string) (*PayloadPaseto, error) {
	panic("not implemented") // TODO: Implement
}
