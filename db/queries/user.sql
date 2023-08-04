-- name: CreateUser :one
INSERT INTO
  users (username, email, hashed_password, full_name)
VALUES
  ($1, $2, $3, $4) RETURNING *;

-- name: GetUser :one
SELECT
  *
FROM
  users
WHERE
  username = $1
LIMIT 1;

-- name: UpdateUserUsingCaseFirst :one
UPDATE users 
SET 
  hashed_password = CASE WHEN sqlc.arg(set_hashed_password)::bool THEN sqlc.arg(hashed_password) ELSE hashed_password END,
  email = CASE WHEN sqlc.arg(set_email)::bool THEN sqlc.arg(email) ELSE email END,
  full_name = CASE WHEN sqlc.arg(set_full_name)::bool THEN sqlc.arg(full_name) ELSE full_name END
WHERE 
  username = sqlc.arg(username)
RETURNING *;

-- name: UpdateUserUsingCaseSecond :one
UPDATE users 
SET 
  hashed_password = coalesce(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = coalesce(sqlc.narg(password_changed_at), password_changed_at),
  email = coalesce(sqlc.narg(email), email),
  full_name = coalesce(sqlc.narg(full_name), full_name),
  is_email_activated = coalesce(sqlc.narg(is_email_activated), is_email_activated)
WHERE 
  username = sqlc.arg('username')
RETURNING *;