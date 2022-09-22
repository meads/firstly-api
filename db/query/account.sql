
-- name: GetAccount :one
SELECT * FROM account
WHERE id = $1 LIMIT 1;

-- name: GetAccountByUsername :one
SELECT * FROM account
WHERE username = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM account LIMIT $1 OFFSET $2;

-- name: CreateAccount :one
INSERT INTO account (
  username, phrase, salt, created
) VALUES (
  $1, $2, $3, NOW()
)
RETURNING *;

-- name: SoftDeleteAccount :exec
UPDATE account
SET deleted = 1
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM account
WHERE id = $1;

-- name: UpdateAccount :exec
UPDATE account
SET username = $1, phrase = $2, salt = $3, updated = NOW()
WHERE id = $4
RETURNING updated;

