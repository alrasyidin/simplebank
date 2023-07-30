package db

import (
	"context"
)

type CreateUserTxParam struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParam) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		user, err := q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		result.User = user

		return arg.AfterCreate(user)
	})
	return result, err
}
