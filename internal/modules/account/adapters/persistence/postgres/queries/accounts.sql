-- name: CreateAccount :exec
INSERT INTO accounts (
  id,
  name,
  balance,
  user_id,
  created_at,
  updated_at,
  deleted_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
;

-- name: ListAccountsByUserID :many
SELECT
  id,
  name,
  balance,
  user_id,
  created_at,
  updated_at,
  deleted_at
FROM accounts
WHERE user_id = $1
  AND deleted_at IS NULL
  AND ($2::timestamptz IS NULL OR (created_at, id) < ($2, $3::uuid))
ORDER BY created_at DESC, id DESC
LIMIT $4
;

-- name: GetAccountByID :one
SELECT
  id,
  name,
  balance,
  user_id,
  created_at,
  updated_at,
  deleted_at
FROM accounts
WHERE id = $1
AND deleted_at IS NULL
;

-- name: UpdateAccountName :exec
UPDATE accounts
SET
  name = $2,
  updated_at = $3
WHERE id = $1
AND deleted_at IS NULL
;

-- name: UpdateAccountBalance :exec
UPDATE accounts
SET
  balance = $2,
  updated_at = $3
WHERE id = $1
AND deleted_at IS NULL
;

-- name: DeleteAccount :exec
UPDATE accounts
SET
  deleted_at = $2,
  updated_at = $3
WHERE id = $1
AND deleted_at IS NULL
;