package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	n            = 10
	amount int64 = 10
)

func TestStore_TransferTx(t *testing.T) {
	store := NewStore(testDB)
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)
	errChan, resultChan := make(chan error), make(chan TransfersTxResult)
	//创建5个携程执行交易操作
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: a1.ID,
				ToAccountID:   a2.ID,
				Amount:        amount,
			})
			errChan <- err
			resultChan <- result
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)
		result := <-resultChan
		require.NotEmpty(t, result)
		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, a1.ID, transfer.FromAccountID)
		require.Equal(t, a2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.CreateAt)
		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, a1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreateAt)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, a2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreateAt)
		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, a1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, a2.ID)

		//check balance
		diff1 := a1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - a2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //amount 2*amount 3*amount ....

		k := int(diff1 / amount)
		require.True(t, k >= 0 && k <= n && k == (i+1))
	}
	//check the final result balance
	updateAccount1, err := testQuery.GetAccount(context.Background(), a1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQuery.GetAccount(context.Background(), a2.ID)
	require.NoError(t, err)

	require.Equal(t, a1.Balance-n*amount, updateAccount1.Balance)
	require.Equal(t, a2.Balance+n*amount, updateAccount2.Balance)
}

func TestStore_TransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	a1 := createRandomAccount(t)
	a2 := createRandomAccount(t)
	errChan := make(chan error)
	//创建5个携程执行交易操作
	for i := 0; i < n; i++ {
		fromID, toID := a1.ID, a2.ID
		if i%2 == 1 {
			fromID, toID = a2.ID, a1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromID,
				ToAccountID:   toID,
				Amount:        amount,
			})
			errChan <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)
	}
	//check the final result balance
	updateAccount1, err := testQuery.GetAccount(context.Background(), a1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQuery.GetAccount(context.Background(), a2.ID)
	require.NoError(t, err)

	require.Equal(t, a1.Balance, updateAccount1.Balance)
	require.Equal(t, a2.Balance, updateAccount2.Balance)
}
