package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	ExecTx(ctx context.Context, fn func(*Queries) error) error
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
	AddMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error)
	Querier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		rbErr := tx.Rollback()

		if err != nil {
			return fmt.Errorf("tx err: %v, rollback error: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParam struct {
	Amount        int64 `json:"amount"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	ToAccount   Account  `json:"to_account"`
	FromAccount Account  `json:"from_account"`
	ToEntry     Entry    `json:"to_entry"`
	FromEntry   Entry    `json:"from_entry"`
}

var txKey = struct{}{}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error
		tx := ctx.Value(txKey)

		fmt.Printf("%s create transfer\n", tx)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Printf("%s from entry\n", tx)
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Printf("%s to entry\n", tx)
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			fmt.Printf("%s update money from account > to account\n", tx)
			result.FromAccount, result.ToAccount, err = store.AddMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			fmt.Printf("%s update money to account > from account\n", tx)
			result.ToAccount, result.FromAccount, err = store.AddMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

func (store *SQLStore) AddMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
