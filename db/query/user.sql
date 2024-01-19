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

-- name: UpdateUser :one
UPDATE users
SET address = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;
