package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {

	testAccount1 := createRandomAccount(t)
	testAccount2 := createRandomAccount(t)

	// run n concurrent transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: testAccount1.ID,
				ToAccountID:   testAccount2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check result

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)

		result := <-results
		assert.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		assert.NotEmpty(t, result)
		assert.Equal(t, testAccount1.ID, transfer.FromAccountID)
		assert.Equal(t, testAccount2.ID, transfer.ToAccountID)
		assert.Equal(t, amount, transfer.Amount)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		assert.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		assert.NotEmpty(t, fromEntry)
		assert.Equal(t, testAccount1.ID, fromEntry.AccountID)
		assert.Equal(t, -amount, fromEntry.Amount)
		assert.NotZero(t, fromEntry.ID)
		assert.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		assert.NoError(t, err)

		toEntry := result.ToEntry
		assert.NotEmpty(t, toEntry)
		assert.Equal(t, testAccount2.ID, toEntry.AccountID)
		assert.Equal(t, amount, toEntry.Amount)
		assert.NotZero(t, toEntry.ID)
		assert.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		assert.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		assert.NotEmpty(t, fromAccount)
		assert.Equal(t, testAccount1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		assert.NotEmpty(t, toAccount)
		assert.Equal(t, testAccount2.ID, toAccount.ID)

		// check account balance
		diff1 := testAccount1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - testAccount2.Balance
		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		assert.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount ... n * amount

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n)
		assert.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final updated balances
	updateAccount1, err := testStore.GetAccount(context.Background(), testAccount1.ID)
	assert.NoError(t, err)

	updateAccount2, err := testStore.GetAccount(context.Background(), testAccount2.ID)
	assert.NoError(t, err)

	assert.Equal(t, testAccount1.Balance-int64(n)*amount, updateAccount1.Balance)
	assert.Equal(t, testAccount2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	testAccount1 := createRandomAccount(t)
	testAccount2 := createRandomAccount(t)

	// run n concurrent transactions
	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := testAccount1.ID
		toAccountID := testAccount2.ID

		if i%2 == 1 {
			fromAccountID = testAccount2.ID
			toAccountID = testAccount1.ID
		}

		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// check result
	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)
	}

	// check final updated balances
	updateAccount1, err := testStore.GetAccount(context.Background(), testAccount1.ID)
	assert.NoError(t, err)

	updateAccount2, err := testStore.GetAccount(context.Background(), testAccount2.ID)
	assert.NoError(t, err)

	assert.Equal(t, testAccount1.Balance, updateAccount1.Balance)
	assert.Equal(t, testAccount2.Balance, updateAccount2.Balance)

}
