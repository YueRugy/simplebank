package db

import (
	"context"
	"database/sql"
	"fmt"
)


type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if errRe := tx.Rollback(); errRe != nil {
			return fmt.Errorf("err : %v,errRe:%v", err, errRe)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransfersTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransfersTxResult, error) {
	var result TransfersTxResult
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		//create transfer
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}
		//add two account entry
		//fromAccount entry
		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}
		//toAccount entry
		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}
		//update accounts balance
		//fromAccount, err := queries.GetAccountForUpdate(ctx, args.FromAccountID)
		//if err != nil {
		//	return err
		//}
		//result.FromAccount, err = queries.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      args.FromAccountID,
		//	Balance: fromAccount.Balance - args.Amount,
		//})
		//update fromAccount balance
		if args.FromAccountID>args.ToAccountID{
			result.FromAccount,err = queries.AddAccountBalance(ctx,AddAccountBalanceParams{
				Amount: -args.Amount,
				ID:     args.FromAccountID,
			})
			//update toAccount balance
			result.ToAccount,err = queries.AddAccountBalance(ctx,AddAccountBalanceParams{
				Amount: args.Amount,
				ID:     args.ToAccountID,
			})
		}else {
			//update toAccount balance
			result.ToAccount,err = queries.AddAccountBalance(ctx,AddAccountBalanceParams{
				Amount: args.Amount,
				ID:     args.ToAccountID,
			})
			result.FromAccount,err = queries.AddAccountBalance(ctx,AddAccountBalanceParams{
				Amount: -args.Amount,
				ID:     args.FromAccountID,
			})
		}

		//toAccount, err := queries.GetAccountForUpdate(ctx, args.ToAccountID)
		//if err != nil {
		//	return err
		//}
		//result.ToAccount, err = queries.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      args.ToAccountID,
		//	Balance: toAccount.Balance + args.Amount,
		//})
		return nil
	})
	return result, err
}
