package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	err = ComparePasswords(password, hashedPassword)
	assert.NoError(t, err)

	wrongPassword := RandomString(6)
	err = ComparePasswords(wrongPassword, hashedPassword)
	assert.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword2)
	assert.NotEqual(t, hashedPassword, hashedPassword2)
}
