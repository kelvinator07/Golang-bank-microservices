package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kelvinator07/golang-bank-microservices/util"
	"github.com/stretchr/testify/assert"
)

// create users, then use user id to create aacount
func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		UserID:        user.ID,
		AccountNumber: util.RandomAccountNumber(),
		Status:        util.RandomStatus(),
		Balance:       util.RandomMoney(),
		CurrencyCode:  util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.CurrencyCode, account.CurrencyCode)

	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	testAccount := createRandomAccount(t)
	expectedAccount, err := testQueries.GetAccount(context.Background(), testAccount.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, expectedAccount)

	assert.Equal(t, testAccount.ID, expectedAccount.ID)
	assert.Equal(t, testAccount.Balance, expectedAccount.Balance)
	assert.Equal(t, testAccount.CurrencyCode, expectedAccount.CurrencyCode)
	assert.WithinDuration(t, testAccount.CreatedAt, expectedAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	testAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      testAccount.ID,
		Balance: util.RandomMoney(),
	}

	expectedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, expectedAccount)

	assert.Equal(t, testAccount.ID, expectedAccount.ID)
	assert.Equal(t, arg.Balance, expectedAccount.Balance)
	assert.Equal(t, testAccount.CurrencyCode, expectedAccount.CurrencyCode)
	assert.WithinDuration(t, testAccount.CreatedAt, expectedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	testAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), testAccount.ID)
	assert.NoError(t, err)

	expectedAccount, err := testQueries.GetAccount(context.Background(), testAccount.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, expectedAccount)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}
	arg := ListAccountsParams{
		UserID: lastAccount.UserID,
		Limit:  5,
		Offset: 0,
	}

	expectedAccounts, err := testQueries.ListAccounts(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, expectedAccounts)

	for _, account := range expectedAccounts {
		assert.NotEmpty(t, account)
		assert.Equal(t, lastAccount.UserID, account.UserID)
	}
}
