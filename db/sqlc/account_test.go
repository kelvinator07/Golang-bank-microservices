package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kelvinator07/golang-bank-microservices/util"
	"github.com/stretchr/testify/assert"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)
	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)

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
	assert.Equal(t, testAccount.Owner, expectedAccount.Owner)
	assert.Equal(t, testAccount.Balance, expectedAccount.Balance)
	assert.Equal(t, testAccount.Currency, expectedAccount.Currency)
	assert.WithinDuration(t, testAccount.CreatedAt, expectedAccount.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	testAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID: testAccount.ID,
		Balance: util.RandomMoney(),
	}

	expectedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, expectedAccount)

	assert.Equal(t, testAccount.ID, expectedAccount.ID)
	assert.Equal(t, testAccount.Owner, expectedAccount.Owner)
	assert.Equal(t, arg.Balance, expectedAccount.Balance)
	assert.Equal(t, testAccount.Currency, expectedAccount.Currency)
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
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}

	expectedAccounts, err := testQueries.ListAccounts(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, expectedAccounts, 5)

	for _, account := range expectedAccounts {
		assert.NotEmpty(t, account)
	}
}
