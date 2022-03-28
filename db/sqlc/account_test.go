package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"
)

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)
	arg := CreateAccountParams{Owner: user.Username, Balance: util.RandomBalance(), Currency: util.RandomCurrency()}
	account, err := testQuery.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreateAt)
}
func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{Owner: user.Username, Balance: util.RandomBalance(), Currency: util.RandomCurrency()}
	account, err := testQuery.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreateAt)
	return account
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreateAt, account2.CreateAt, 1*time.Second)
}

func TestQueries_UpdateAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)
	param := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomBalance(),
	}
	account2, err := testQuery.UpdateAccount(context.Background(), param)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, param.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreateAt, account2.CreateAt, 1*time.Second)
}

func TestQueries_DeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQuery.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	account2, err := testQuery.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestQueries_ListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 5; i++ {
		lastAccount = createRandomAccount(t)
	}
	args := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}
	list, err := testQuery.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, list)
	for _, temp := range list {
		require.NotEmpty(t, temp)
		require.Equal(t, temp.Owner, lastAccount.Owner)
	}
}

func TestQueries_AddAccountBalance(t *testing.T) {
	account := createRandomAccount(t)
	updatedAccount, err := testQuery.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		Amount: amount,
		ID:     account.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, account.ID, updatedAccount.ID)
	require.Equal(t, account.Balance+amount, updatedAccount.Balance)
	require.WithinDuration(t, account.CreateAt, updatedAccount.CreateAt, 1*time.Second)
}

func TestQueries_GetAccountForUpdate(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQuery.GetAccountForUpdate(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreateAt, account2.CreateAt, 1*time.Second)
}
