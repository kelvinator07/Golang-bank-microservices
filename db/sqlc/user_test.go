package db

import (
	"context"
	"testing"

	"github.com/kelvinator07/golang-bank-microservices/util"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

// create users, then use user id to create aacount
func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	assert.NoError(t, err)

	arg := CreateUserParams{
		AccountName:    util.RandomString(10),
		HashedPassword: hashedPassword,
		Address:        util.RandomString(20),
		Gender:         util.RandomGender(),
		PhoneNumber:    util.RandomPhoneNumber(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.NotZero(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
	assert.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestListUser(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	expectedUsers, err := testStore.ListUsers(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, expectedUsers, 5)

	for _, user := range expectedUsers {
		assert.NotEmpty(t, user)
	}
}
