-- name: CreateUser :one
INSERT INTO users (
  account_name,
  hashed_password,
  address,
  gender,
  phone_number,
  email
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1 
FOR NO KEY UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetAllUsers :many
SELECT * FROM users
WHERE id > $1
ORDER BY id
LIMIT $2;

-- name: UpdateUser :one
UPDATE users
SET
  account_name = COALESCE(sqlc.narg(account_name), account_name),
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  address = COALESCE(sqlc.narg(address), address),
  phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
  email = COALESCE(sqlc.narg(email), email),
  is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified)
WHERE
  email = sqlc.arg(email)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;
