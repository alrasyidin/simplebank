package token

import "time"

type Generator interface {
	CreateToken(username string, expire time.Duration) (string, *PayloadJWT, error)
	VerifyToken(token string) (*PayloadJWT, error)
	VerifyTokenPaseto(token string) (*PayloadPaseto, error)
}
