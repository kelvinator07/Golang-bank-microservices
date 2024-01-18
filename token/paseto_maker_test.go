package token

import (
	"testing"
	"time"

	"github.com/kelvinator07/golang-bank-microservices/util"
	"github.com/stretchr/testify/assert"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	assert.NoError(t, err)

	accountName := util.RandomAccountName()
	email := util.RandomEmail()

	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(accountName, email, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)

	assert.NotZero(t, payload.ID)
	assert.Equal(t, accountName, payload.AccountName)
	assert.Equal(t, email, payload.Email)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	assert.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoMakerToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	assert.NoError(t, err)

	accountName := util.RandomAccountName()
	email := util.RandomEmail()

	expiredDuration := -time.Minute

	token, err := maker.CreateToken(accountName, email, expiredDuration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrExpiredToken.Error())

	assert.Nil(t, payload)
}

func TestInvalidPasetoMakerToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	assert.NoError(t, err)

	accountName := util.RandomAccountName()
	email := util.RandomEmail()

	duration := time.Minute

	token, err := maker.CreateToken(accountName, email, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	token = token[:len(token)-6] + "invaliddata"

	payload, err := maker.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrInvalidToken.Error())

	assert.Nil(t, payload)
}
