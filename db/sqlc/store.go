package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	ExecTx(ctx context.Context, fn func(*Queries) error) error
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParam) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParam) (VerifyEmailTxResult, error)
	AddMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error)
	Querier
}

type SQLStore struct {
	*Queries
	poolConn *pgxpool.Pool
}

func NewStore(poolConn *pgxpool.Pool) Store {
	return &SQLStore{
		Queries:  New(poolConn),
		poolConn: poolConn,
	}
}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.poolConn.Begin(ctx)

	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		rbErr := tx.Rollback(ctx)

		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback error: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit(ctx)
}
