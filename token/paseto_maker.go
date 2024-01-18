package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

const symmetricKeySize = chacha20poly1305.KeySize

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != symmetricKeySize {
		return nil, fmt.Errorf("invalid key size, must be exactly %d characters", symmetricKeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (pm *PasetoMaker) CreateToken(accountName string, email string, duration time.Duration) (string, error) {
	payload, err := NewPayload(accountName, email, duration)
	if err != nil {
		return "", err
	}
	return pm.paseto.Encrypt(pm.symmetricKey, payload, nil)
}

func (pm *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := pm.paseto.Decrypt(token, pm.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
