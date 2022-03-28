package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"
)

func TestQueries_CreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomAmount(account.Balance),
	}
	entry, err := testQuery.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreateAt)
}
func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomAmount(account.Balance),
	}
	entry, err := testQuery.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreateAt)

	return entry
}

func TestQueries_GetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	entry2, err := testQuery.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreateAt, entry2.CreateAt, 1*time.Second)
}

func TestQueries_UpdateEntryAmount(t *testing.T) {
	entry1 := createRandomEntry(t)
	args := UpdateEntryAmountParams{
		ID:     entry1.ID,
		Amount: util.RandomAmount(entry1.Amount),
	}
	entry2, err := testQuery.UpdateEntryAmount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, args.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreateAt, entry2.CreateAt, 1*time.Second)
}

func TestQueries_DeleteEntry(t *testing.T) {
	entry1 := createRandomEntry(t)
	err := testQuery.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	entry2, err := testQuery.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestQueries_ListEntries(t *testing.T) {
	for i := 0; i < 5; i++ {
		_ = createRandomEntry(t)
	}
	args := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}
	list, err := testQuery.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, list, 5)
	for _, entry := range list {
		require.NotEmpty(t, entry)
	}
}
