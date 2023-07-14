package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)
	n := 5
	amount := int64(10)

	errs := make(chan error)

	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParam{
				Amount:        amount,
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
			})

			errs <- err
			results <- result
		}()
	}
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)

		// check transfer object
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		_, err = store.Queries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println(">> tx: ", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		fmt.Println("diff", diff1, diff2)
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before: ", account1.Balance, account2.Balance)
	n := 10
	amount := int64(10)

	errs := make(chan error)
	// results := make(chan TransferTxResult)
	// var mutex sync.Mutex
	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}
		go func() {
			// mutex.Lock()
			_, err := store.TransferTx(context.Background(), TransferTxParam{
				Amount:        amount,
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
			})
			errs <- err
			// mutex.Unlock()
			// results <- result
		}()
	}

	// var existed = make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	// check final update balance
	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
	fmt.Println(">> after: ", updateAccount1.Balance, updateAccount2.Balance)
	fmt.Println(">> test: ", account1, account2)
}
