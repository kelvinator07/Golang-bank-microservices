-- name: CreateVerifyEmail :one
INSERT INTO verify_emails (
  email,
  secret_code
) VALUES (
  $1, $2
) RETURNING *;

-- name: UpdateVerifyEmail :one
UPDATE verify_emails
SET is_used = TRUE
WHERE id = $1
RETURNING *;

-- name: GetVerifyEmail :one
SELECT * FROM verify_emails
WHERE id = $1 LIMIT 1;


