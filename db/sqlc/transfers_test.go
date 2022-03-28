package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"
)

func TestQueries_CreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomAmount(account1.Balance),
	}
	transfer, err := testQuery.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreateAt)
}

func getRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomAmount(account1.Balance),
	}
	transfer, err := testQuery.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreateAt)
	return transfer
}

func TestQueries_GetTransfer(t *testing.T) {
	transfer1 := getRandomTransfer(t)
	transfer2, err := testQuery.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreateAt, transfer2.CreateAt, 1*time.Second)
}

func TestQueries_UpdateTransferAmount(t *testing.T) {
	transfer1 :=getRandomTransfer(t)
	args := UpdateTransferAmountParams{
		ID:     transfer1.ID,
		Amount: util.RandomAmount(transfer1.Amount),
	}
	transfer2, err := testQuery.UpdateTransferAmount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, args.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreateAt, transfer2.CreateAt, 1*time.Second)
}

func TestQueries_DeleteTransfer(t *testing.T) {
	transfer1 := getRandomTransfer(t)
	err := testQuery.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	transfer2, err := testQuery.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestQueries_ListTransfers(t *testing.T) {
	for i := 0; i < 5; i++ {
		_ = getRandomTransfer(t)
	}
	args := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	list, err := testQuery.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, list, 5)
	for _, transfer := range list {
		require.NotEmpty(t, transfer)
	}
}
