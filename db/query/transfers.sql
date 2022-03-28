-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id,
                       to_account_id,
                     amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTransfer :one
SELECT *
FROM transfers
WHERE id = $1;

-- name: ListTransfers :many
SELECT *
FROM transfers
ORDER BY id
limit $1 offset $2
;

-- name: UpdateTransferAmount :one
UPDATE transfers
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTransfer :exec
DELETE
FROM transfers
WHERE id = $1;
