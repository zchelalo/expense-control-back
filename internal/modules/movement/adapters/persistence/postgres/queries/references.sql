-- name: UserExists :one
SELECT EXISTS (
  SELECT 1
  FROM users
  WHERE id = $1
    AND deleted_at IS NULL
)
;

-- name: AccountExistsByUserID :one
SELECT EXISTS (
  SELECT 1
  FROM accounts
  WHERE id = $1
    AND user_id = $2
    AND deleted_at IS NULL
)
;

-- name: IncreaseAccountBalance :execrows
UPDATE accounts
SET
  balance = balance + $1,
  updated_at = $4
WHERE id = $2
  AND user_id = $3
  AND deleted_at IS NULL
;

-- name: DecreaseAccountBalance :execrows
UPDATE accounts
SET
  balance = balance - $1,
  updated_at = $4
WHERE id = $2
  AND user_id = $3
  AND deleted_at IS NULL
  AND balance >= $1
;

-- name: GetMovementTypeByID :one
SELECT
  id,
  name,
  key,
  description,
  created_at,
  updated_at,
  deleted_at
FROM movement_types
WHERE id = $1
  AND deleted_at IS NULL
;

-- name: GetCategoryByID :one
SELECT
  id,
  name,
  created_at,
  updated_at,
  deleted_at
FROM categories
WHERE id = $1
  AND deleted_at IS NULL
;
