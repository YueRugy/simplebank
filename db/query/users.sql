-- name: CreateUser :one
INSERT INTO users (username,
                   hashed_password,
                   email,
                   full_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUsers :one
SELECT *
FROM users
WHERE username = $1
limit 1;

/*-- name: GetAccountForUpdate :one
SELECT *
FROM account
WHERE id = $1
FOR no key update ;*/


/*-- name: ListUsers :many
SELECT *
FROM users
ORDER BY username
limit $1 offset $2
;*/

/*-- name: UpdateUsers :one
UPDATE users
SET hashed_password = $2
WHERE username = $1
RETURNING *;*/
/*-- name: AddAccountBalance :one
UPDATE account
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;*/

/*-- name: DeleteUser :exec
DELETE
FROM users
WHERE username = $1;*/
