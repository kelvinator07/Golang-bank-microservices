package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.GetVerifyEmail(ctx, arg.EmailId)
		if err != nil {
			if err == sql.ErrNoRows {
				err = fmt.Errorf("verify email with id %v doesn't exist", arg.EmailId)
				return err
			}
			return err
		}

		if result.VerifyEmail.SecretCode != arg.SecretCode {
			err = fmt.Errorf("verify code invalid")
			return err
		}

		if result.VerifyEmail.IsUsed {
			err = fmt.Errorf("verify code already used")
			return err
		}

		if time.Now().After(result.VerifyEmail.ExpiredAt) {
			err = fmt.Errorf("verify code already expired")
			return err
		}

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, arg.EmailId)
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Email: sql.NullString{
				String: result.VerifyEmail.Email,
				Valid:  true,
			},
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})

		return err
	})

	return result, err
}
