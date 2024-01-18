package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token authentication")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload data of the token
type Payload struct {
	ID          uuid.UUID `json:"id"`
	AccountName string    `json:"account_name"`
	Email       string    `json:"email"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func NewPayload(accountName string, email string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom() // len 32
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:          tokenID,
		AccountName: accountName,
		Email:       email,
		IssuedAt:    time.Now(),
		ExpiredAt:   time.Now().Add(duration),
	}
	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) { // current time more than expired at
		return ErrExpiredToken
	}
	return nil
}
