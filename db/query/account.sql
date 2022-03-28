-- name: CreateAccount :one
INSERT INTO account (owner,
                     balance,
                     currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccount :one
SELECT *
FROM account
WHERE id = $1;

-- name: GetAccountForUpdate :one
SELECT *
FROM account
WHERE id = $1
FOR no key update ;


-- name: ListAccounts :many
SELECT *
FROM account
where owner=$1
ORDER BY id
limit $2 offset $3
;

-- name: UpdateAccount :one
UPDATE account
SET balance = $2
WHERE id = $1
RETURNING *;
-- name: AddAccountBalance :one
UPDATE account
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE
FROM account
WHERE id = $1;
