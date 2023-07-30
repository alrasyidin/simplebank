package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	ExecTx(ctx context.Context, fn func(*Queries) error) error
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParam) (CreateUserTxResult, error)
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

		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback error: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}
