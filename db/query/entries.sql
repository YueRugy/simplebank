-- name: CreateEntry :one
INSERT INTO entries (account_id,
                     amount)
VALUES ($1, $2)
RETURNING *;

-- name: GetEntry :one
SELECT *
FROM entries
WHERE id = $1;

-- name: ListEntries :many
SELECT *
FROM entries
ORDER BY id
limit $1 offset $2
;

-- name: UpdateEntryAmount :one
UPDATE entries
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE
FROM entries
WHERE id = $1;
