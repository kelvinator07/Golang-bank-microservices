// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: verify_email.sql

package db

import (
	"context"
)

const createVerifyEmail = `-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  email,
  secret_code
) VALUES (
  $1, $2
) RETURNING id, email, secret_code, is_used, created_at, expired_at
`

type CreateVerifyEmailParams struct {
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (VerifyEmail, error) {
	row := q.db.QueryRow(ctx, createVerifyEmail, arg.Email, arg.SecretCode)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const getVerifyEmail = `-- name: GetVerifyEmail :one
SELECT id, email, secret_code, is_used, created_at, expired_at FROM verify_emails
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetVerifyEmail(ctx context.Context, id int64) (VerifyEmail, error) {
	row := q.db.QueryRow(ctx, getVerifyEmail, id)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}

const updateVerifyEmail = `-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used = TRUE
WHERE id = $1
RETURNING id, email, secret_code, is_used, created_at, expired_at
`

func (q *Queries) UpdateVerifyEmail(ctx context.Context, id int64) (VerifyEmail, error) {
	row := q.db.QueryRow(ctx, updateVerifyEmail, id)
	var i VerifyEmail
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.SecretCode,
		&i.IsUsed,
		&i.CreatedAt,
		&i.ExpiredAt,
	)
	return i, err
}
