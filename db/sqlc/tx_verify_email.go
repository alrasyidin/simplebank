package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailTxParam struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParam) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult
	err := store.ExecTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})

		if err != nil {
			return err
		}

		result.User, err = q.UpdateUserUsingCaseSecond(ctx, UpdateUserUsingCaseSecondParams{
			Username: result.VerifyEmail.Username,
			IsEmailActivated: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})

		if err != nil {
			return err
		}

		return err
	})
	return result, err
}
