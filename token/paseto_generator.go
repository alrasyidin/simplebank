package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoGenerator struct {
	paseto      paseto.V2
	symetricKey string
}

const minSymetricKeySize = 32

// Constructor for PasetoGenerator
func NewPasetoGenerator(symetricKey string) (*PasetoGenerator, error) {
	if len(symetricKey) < minSymetricKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSymetricKeySize)
	}

	return &PasetoGenerator{
		paseto:      *paseto.NewV2(),
		symetricKey: symetricKey,
	}, nil
}

func (generator *PasetoGenerator) CreateToken(username string, expire time.Duration) (string, error) {
	payload, err := NewPayloadPaseto(username, expire)

	if err != nil {
		return "", err
	}

	return generator.paseto.Encrypt([]byte(generator.symetricKey), payload, nil)
}

func (generator *PasetoGenerator) VerifyTokenPaseto(token string) (*PayloadPaseto, error) {
	payload := &PayloadPaseto{}

	err := generator.paseto.Decrypt(token, []byte(generator.symetricKey), payload, nil)

	if err != nil {
		return nil, ErrInvalidToken
	}
	if err = payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
