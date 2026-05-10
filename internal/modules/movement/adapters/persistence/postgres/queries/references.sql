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

-- name: ListMovementTypes :many
SELECT
  id,
  name,
  key,
  description,
  created_at,
  updated_at,
  deleted_at
FROM movement_types
WHERE deleted_at IS NULL
ORDER BY key ASC, id ASC
;

-- name: GetCategoryByIDForUser :one
SELECT
  c.id,
  c.name,
  c.is_system,
  c.system_key
FROM user_categories uc
INNER JOIN categories c
  ON c.id = uc.category_id
 AND c.deleted_at IS NULL
WHERE uc.category_id = sqlc.arg('category_id')
  AND uc.user_id = sqlc.arg('user_id')
  AND uc.deleted_at IS NULL
  AND c.is_system = FALSE
;
