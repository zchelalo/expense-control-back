-- name: GetMovementTypeIDByKey :one
SELECT id
FROM movement_types
WHERE key = sqlc.arg('key')
  AND deleted_at IS NULL
;

-- name: UpsertSystemCategoryByKey :one
INSERT INTO categories (
  id,
  name,
  is_system,
  system_key,
  created_at,
  updated_at,
  deleted_at
)
VALUES (
  sqlc.arg('id'),
  sqlc.arg('name'),
  TRUE,
  sqlc.arg('system_key'),
  sqlc.arg('created_at'),
  sqlc.arg('updated_at'),
  sqlc.arg('deleted_at')
)
ON CONFLICT (system_key) WHERE system_key IS NOT NULL DO UPDATE
SET
  name = EXCLUDED.name,
  deleted_at = NULL,
  updated_at = EXCLUDED.updated_at
RETURNING id
;

-- name: UpsertUserCategory :one
INSERT INTO user_categories (
  user_id,
  category_id,
  created_at,
  updated_at,
  deleted_at
)
VALUES (
  sqlc.arg('user_id'),
  sqlc.arg('category_id'),
  sqlc.arg('created_at'),
  sqlc.arg('updated_at'),
  sqlc.arg('deleted_at')
)
ON CONFLICT (user_id, category_id) DO UPDATE
SET
  updated_at = CASE
    WHEN user_categories.deleted_at IS NULL THEN user_categories.updated_at
    ELSE EXCLUDED.updated_at
  END,
  deleted_at = NULL
RETURNING user_id, category_id, created_at, updated_at
;

-- name: CreateInitialBalanceMovement :exec
INSERT INTO movements (
  id,
  amount,
  description,
  movement_type_id,
  category_id,
  account_id,
  user_id,
  created_at,
  updated_at,
  deleted_at
)
VALUES (
  sqlc.arg('id'),
  sqlc.arg('amount'),
  sqlc.arg('description'),
  sqlc.arg('movement_type_id'),
  sqlc.arg('category_id'),
  sqlc.arg('account_id'),
  sqlc.arg('user_id'),
  sqlc.arg('created_at'),
  sqlc.arg('updated_at'),
  sqlc.arg('deleted_at')
)
;
