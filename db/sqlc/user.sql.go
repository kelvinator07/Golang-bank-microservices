// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  account_name,
  hashed_password,
  address,
  gender,
  phone_number,
  email
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING id, account_name, hashed_password, address, gender, phone_number, email, password_changed_at, created_at
`

type CreateUserParams struct {
	AccountName    string `json:"account_name"`
	HashedPassword string `json:"hashed_password"`
	Address        string `json:"address"`
	Gender         string `json:"gender"`
	PhoneNumber    int64  `json:"phone_number"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.AccountName,
		arg.HashedPassword,
		arg.Address,
		arg.Gender,
		arg.PhoneNumber,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.AccountName,
		&i.HashedPassword,
		&i.Address,
		&i.Gender,
		&i.PhoneNumber,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, account_name, hashed_password, address, gender, phone_number, email, password_changed_at, created_at FROM users
WHERE id > $1
ORDER BY id
LIMIT $2
`

type GetAllUsersParams struct {
	ID    int64 `json:"id"`
	Limit int32 `json:"limit"`
}

func (q *Queries) GetAllUsers(ctx context.Context, arg GetAllUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.AccountName,
			&i.HashedPassword,
			&i.Address,
			&i.Gender,
			&i.PhoneNumber,
			&i.Email,
			&i.PasswordChangedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT id, account_name, hashed_password, address, gender, phone_number, email, password_changed_at, created_at FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.AccountName,
		&i.HashedPassword,
		&i.Address,
		&i.Gender,
		&i.PhoneNumber,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, account_name, hashed_password, address, gender, phone_number, email, password_changed_at, created_at FROM users
WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.AccountName,
		&i.HashedPassword,
		&i.Address,
		&i.Gender,
		&i.PhoneNumber,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getUserForUpdate = `-- name: GetUserForUpdate :one
SELECT id, account_name, hashed_password, address, gender, phone_number, email, password_changed_at, created_at FROM users
WHERE id = $1 LIMIT 1 
FOR NO KEY UPDATE
`

func (q *Queries) GetUserForUpdate(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserForUpdate, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.AccountName,
		&i.HashedPassword,
		&i.Address,
		&i.Gender,
		&i.PhoneNumber,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, account_name, hashed_password, address, gender, phone_number, email, password_changed_at, created_at FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.AccountName,
			&i.HashedPassword,
			&i.Address,
			&i.Gender,
			&i.PhoneNumber,
			&i.Email,
			&i.PasswordChangedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET address = $2
WHERE id = $1
RETURNING id, account_name, hashed_password, address, gender, phone_number, email, password_changed_at, created_at
`

type UpdateUserParams struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.ID, arg.Address)
	var i User
	err := row.Scan(
		&i.ID,
		&i.AccountName,
		&i.HashedPassword,
		&i.Address,
		&i.Gender,
		&i.PhoneNumber,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
