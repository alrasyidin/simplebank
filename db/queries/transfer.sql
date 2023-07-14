-- name: CreateTransfer :one
INSERT INTO
  transfers (from_account_id, to_account_id, amount)
VALUES
  ($1, $2, $3) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers WHERE id = $1 LIMIT 1;

-- name: ListTransfer :many
SELECT * FROM transfers;

-- name: UpdateTransfer :one
UPDATE transfers SET amount = $2 WHERE id = $1 RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfers WHERE id = $1;