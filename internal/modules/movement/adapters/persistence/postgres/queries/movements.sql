-- name: CreateMovement :exec
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
;

-- name: GetMovementByID :one
SELECT
  m.id,
  m.amount,
  m.description,
  m.movement_type_id,
  m.category_id,
  m.account_id,
  m.user_id,
  m.created_at,
  m.updated_at,
  m.deleted_at
FROM movements m
WHERE m.id = $1
  AND m.deleted_at IS NULL
;

-- name: GetMovementDetailsByIDForUser :one
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.id = $1
    AND m.user_id = $2
    AND m.deleted_at IS NULL
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
;

-- name: ListMovementsByUserIDAfter :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.user_id = sqlc.arg('user_id')
    AND m.deleted_at IS NULL
    AND (
      sqlc.narg('cursor_created_at')::timestamptz IS NULL
      OR (m.created_at, m.id) < (
        sqlc.narg('cursor_created_at')::timestamptz,
        sqlc.narg('cursor_movement_id')::uuid
      )
    )
  ORDER BY m.created_at DESC, m.id DESC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at DESC, base.id DESC
;

-- name: ListMovementsByUserIDBefore :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.user_id = sqlc.arg('user_id')
    AND m.deleted_at IS NULL
    AND (m.created_at, m.id) > (
      sqlc.arg('cursor_created_at')::timestamptz,
      sqlc.arg('cursor_movement_id')::uuid
    )
  ORDER BY m.created_at ASC, m.id ASC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at ASC, base.id ASC
;

-- name: ListMovementsByUserIDAndAccountIDAfter :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.account_id = sqlc.arg('account_id')
    AND m.user_id = sqlc.arg('user_id')
    AND m.deleted_at IS NULL
    AND (
      sqlc.narg('cursor_created_at')::timestamptz IS NULL
      OR (m.created_at, m.id) < (
        sqlc.narg('cursor_created_at')::timestamptz,
        sqlc.narg('cursor_movement_id')::uuid
      )
    )
  ORDER BY m.created_at DESC, m.id DESC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at DESC, base.id DESC
;

-- name: ListMovementsByUserIDAndAccountIDBefore :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.account_id = sqlc.arg('account_id')
    AND m.user_id = sqlc.arg('user_id')
    AND m.deleted_at IS NULL
    AND (m.created_at, m.id) > (
      sqlc.arg('cursor_created_at')::timestamptz,
      sqlc.arg('cursor_movement_id')::uuid
    )
  ORDER BY m.created_at ASC, m.id ASC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at ASC, base.id ASC
;

-- name: ListMovementsByUserIDAndCategoryIDAfter :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.user_id = sqlc.arg('user_id')
    AND m.category_id = sqlc.arg('category_id')
    AND m.deleted_at IS NULL
    AND (
      sqlc.narg('cursor_created_at')::timestamptz IS NULL
      OR (m.created_at, m.id) < (
        sqlc.narg('cursor_created_at')::timestamptz,
        sqlc.narg('cursor_movement_id')::uuid
      )
    )
  ORDER BY m.created_at DESC, m.id DESC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at DESC, base.id DESC
;

-- name: ListMovementsByUserIDAndCategoryIDBefore :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.user_id = sqlc.arg('user_id')
    AND m.category_id = sqlc.arg('category_id')
    AND m.deleted_at IS NULL
    AND (m.created_at, m.id) > (
      sqlc.arg('cursor_created_at')::timestamptz,
      sqlc.arg('cursor_movement_id')::uuid
    )
  ORDER BY m.created_at ASC, m.id ASC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at ASC, base.id ASC
;

-- name: ListMovementsByUserIDAndMovementTypeIDAfter :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.user_id = sqlc.arg('user_id')
    AND m.movement_type_id = sqlc.arg('movement_type_id')
    AND m.deleted_at IS NULL
    AND (
      sqlc.narg('cursor_created_at')::timestamptz IS NULL
      OR (m.created_at, m.id) < (
        sqlc.narg('cursor_created_at')::timestamptz,
        sqlc.narg('cursor_movement_id')::uuid
      )
    )
  ORDER BY m.created_at DESC, m.id DESC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at DESC, base.id DESC
;

-- name: ListMovementsByUserIDAndMovementTypeIDBefore :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.user_id = sqlc.arg('user_id')
    AND m.movement_type_id = sqlc.arg('movement_type_id')
    AND m.deleted_at IS NULL
    AND (m.created_at, m.id) > (
      sqlc.arg('cursor_created_at')::timestamptz,
      sqlc.arg('cursor_movement_id')::uuid
    )
  ORDER BY m.created_at ASC, m.id ASC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at ASC, base.id ASC
;

-- name: ListMovementsByUserIDFilteredAfter :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.user_id = sqlc.arg('user_id')
    AND m.deleted_at IS NULL
    AND (
      sqlc.narg('account_id')::uuid IS NULL
      OR m.account_id = sqlc.narg('account_id')::uuid
    )
    AND (
      sqlc.narg('category_id')::uuid IS NULL
      OR m.category_id = sqlc.narg('category_id')::uuid
    )
    AND (
      sqlc.narg('movement_type_id')::uuid IS NULL
      OR m.movement_type_id = sqlc.narg('movement_type_id')::uuid
    )
    AND (
      sqlc.narg('cursor_created_at')::timestamptz IS NULL
      OR (m.created_at, m.id) < (
        sqlc.narg('cursor_created_at')::timestamptz,
        sqlc.narg('cursor_movement_id')::uuid
      )
    )
  ORDER BY m.created_at DESC, m.id DESC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at DESC, base.id DESC
;

-- name: ListMovementsByUserIDFilteredBefore :many
WITH base AS (
  SELECT
    m.id,
    m.amount,
    m.description,
    m.movement_type_id,
    m.category_id,
    m.account_id,
    m.user_id,
    m.created_at,
    m.updated_at,
    m.deleted_at
  FROM movements m
  WHERE m.user_id = sqlc.arg('user_id')
    AND m.deleted_at IS NULL
    AND (
      sqlc.narg('account_id')::uuid IS NULL
      OR m.account_id = sqlc.narg('account_id')::uuid
    )
    AND (
      sqlc.narg('category_id')::uuid IS NULL
      OR m.category_id = sqlc.narg('category_id')::uuid
    )
    AND (
      sqlc.narg('movement_type_id')::uuid IS NULL
      OR m.movement_type_id = sqlc.narg('movement_type_id')::uuid
    )
    AND (m.created_at, m.id) > (
      sqlc.arg('cursor_created_at')::timestamptz,
      sqlc.arg('cursor_movement_id')::uuid
    )
  ORDER BY m.created_at ASC, m.id ASC
  LIMIT sqlc.arg('limit_count')
)
SELECT
  base.id,
  base.amount,
  base.description,
  base.movement_type_id,
  base.category_id,
  base.account_id,
  base.user_id,
  base.created_at,
  base.updated_at,
  base.deleted_at,
  mt.key AS movement_type_key,
  mt.name AS movement_type_name,
  c.name AS category_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
ORDER BY base.created_at ASC, base.id ASC
;

-- name: DeleteMovement :execrows
UPDATE movements
SET
  deleted_at = $2,
  updated_at = $3
WHERE id = $1
  AND deleted_at IS NULL
;
