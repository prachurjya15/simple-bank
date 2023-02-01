package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	n := 10
	amount := int64(100)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	var txName string
	for i := 0; i < n; i++ {
		i := i
		go func() {
			txName = fmt.Sprintf("Tx---> %d", i+1)
			ctx := context.WithValue(context.Background(), TxKey, txName)
			res, err := store.TransferTx(ctx, TransferTxParam{account1.ID, account2.ID, amount})
			errs <- err
			results <- res
		}()
	}
	existed := make(map[int]bool)
	log.Printf("Before Transfer  acc1 : %#v, acc2: %#v", account1.Balance, account2.Balance)
	for i := 0; i < n; i++ {
		log.Printf("**********TRANSACTION Number: %d ********************* \n", i+1)
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		require.Equal(t, result.Transfer.FromAccountID, account1.ID)
		require.Equal(t, result.Transfer.ToAccountID, account2.ID)
		require.Equal(t, result.Transfer.Amount, amount)

		_, err = store.Query.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		toEntry := result.ToEntry

		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		_, err = store.Query.GetEntryById(context.Background(), result.FromEntry.ID)
		require.NoError(t, err)

		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		_, err = store.Query.GetEntryById(context.Background(), result.ToEntry.ID)
		require.NoError(t, err)
		// TODO: Update the balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		log.Printf("After Transfer  acc1 : %#v, acc2: %#v", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := account2.Balance - toAccount.Balance
		require.True(t, diff1+diff2 == 0)
		require.True(t, diff1 > 0)
		require.True(t, diff2 < 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)

		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	// Check the updated balance of both account
	updatedAcc1, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAcc1)
	require.NoError(t, err)

	updatedAcc2, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAcc2)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAcc1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAcc2.Balance)

}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	n := 10
	amount := int64(100)

	errs := make(chan error)
	var txName string
	for i := 0; i < n; i++ {
		fromAccount := account1.ID
		toAccount := account2.ID
		i := i
		if i%2 == 1 {
			fromAccount = account2.ID
			toAccount = account1.ID
		}
		go func() {
			txName = fmt.Sprintf("Tx---> %d", i+1)
			ctx := context.WithValue(context.Background(), TxKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParam{fromAccount, toAccount, amount})
			errs <- err
		}()
	}
	log.Printf("Before Transfer  acc1 : %#v, acc2: %#v", account1.Balance, account2.Balance)
	for i := 0; i < n; i++ {
		log.Printf("**********TRANSACTION Number: %d ********************* \n", i+1)
		err := <-errs
		require.NoError(t, err)
	}
	// Check the updated balance of both account
	updatedAcc1, err := testQueries.GetAccountById(context.Background(), account1.ID)
	require.NotEmpty(t, updatedAcc1)
	require.NoError(t, err)

	updatedAcc2, err := testQueries.GetAccountById(context.Background(), account2.ID)
	require.NotEmpty(t, updatedAcc2)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAcc1.Balance)
	require.Equal(t, account2.Balance, updatedAcc2.Balance)

}
