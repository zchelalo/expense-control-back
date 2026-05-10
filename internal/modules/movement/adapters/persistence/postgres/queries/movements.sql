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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
    AND (
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
    AND (
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
    AND (
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
    AND (
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
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
    AND (
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
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
  c.name AS category_name,
  c.is_system AS category_is_system,
  c.system_key AS category_system_key,
  a.name AS account_name
FROM base
INNER JOIN movement_types mt
  ON mt.id = base.movement_type_id
 AND mt.deleted_at IS NULL
INNER JOIN categories c
  ON c.id = base.category_id
 AND c.deleted_at IS NULL
INNER JOIN accounts a
  ON a.id = base.account_id
 AND a.deleted_at IS NULL
ORDER BY base.created_at ASC, base.id ASC
;

-- name: GetMovementStatsOverviewByUserID :one
WITH filtered AS (
  SELECT
    m.amount,
    mt.key AS movement_type_key
  FROM movements m
  INNER JOIN movement_types mt
    ON mt.id = m.movement_type_id
   AND mt.deleted_at IS NULL
  INNER JOIN categories c
    ON c.id = m.category_id
   AND c.deleted_at IS NULL
  INNER JOIN accounts a
    ON a.id = m.account_id
   AND a.deleted_at IS NULL
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
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
    )
)
SELECT
  COUNT(*)::bigint AS total_movements,
  COUNT(*) FILTER (WHERE movement_type_key = 'income')::bigint AS income_count,
  COUNT(*) FILTER (WHERE movement_type_key = 'expense')::bigint AS expense_count,
  COALESCE(SUM(amount) FILTER (WHERE movement_type_key = 'income'), 0)::numeric AS income_total,
  COALESCE(SUM(amount) FILTER (WHERE movement_type_key = 'expense'), 0)::numeric AS expense_total,
  COALESCE(SUM(CASE
    WHEN movement_type_key = 'income' THEN amount
    WHEN movement_type_key = 'expense' THEN -amount
    ELSE 0
  END), 0)::numeric AS net_total
FROM filtered
;

-- name: ListMovementStatsByAccountByUserID :many
WITH filtered AS (
  SELECT
    a.id AS account_id,
    a.name AS account_name,
    m.amount,
    mt.key AS movement_type_key
  FROM movements m
  INNER JOIN movement_types mt
    ON mt.id = m.movement_type_id
   AND mt.deleted_at IS NULL
  INNER JOIN categories c
    ON c.id = m.category_id
   AND c.deleted_at IS NULL
  INNER JOIN accounts a
    ON a.id = m.account_id
   AND a.deleted_at IS NULL
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
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
    )
)
SELECT
  account_id,
  account_name,
  COUNT(*)::bigint AS movement_count,
  COUNT(*) FILTER (WHERE movement_type_key = 'income')::bigint AS income_count,
  COUNT(*) FILTER (WHERE movement_type_key = 'expense')::bigint AS expense_count,
  COALESCE(SUM(amount) FILTER (WHERE movement_type_key = 'income'), 0)::numeric AS income_total,
  COALESCE(SUM(amount) FILTER (WHERE movement_type_key = 'expense'), 0)::numeric AS expense_total,
  COALESCE(SUM(CASE
    WHEN movement_type_key = 'income' THEN amount
    WHEN movement_type_key = 'expense' THEN -amount
    ELSE 0
  END), 0)::numeric AS net_total
FROM filtered
GROUP BY account_id, account_name
ORDER BY expense_total DESC, movement_count DESC, account_name ASC
;

-- name: ListMovementStatsByCategoryByUserID :many
WITH filtered AS (
  SELECT
    c.id AS category_id,
    c.name AS category_name,
    c.is_system AS category_is_system,
    c.system_key AS category_system_key,
    m.amount,
    mt.key AS movement_type_key
  FROM movements m
  INNER JOIN movement_types mt
    ON mt.id = m.movement_type_id
   AND mt.deleted_at IS NULL
  INNER JOIN categories c
    ON c.id = m.category_id
   AND c.deleted_at IS NULL
  INNER JOIN accounts a
    ON a.id = m.account_id
   AND a.deleted_at IS NULL
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
      sqlc.narg('date_from')::timestamptz IS NULL
      OR m.created_at >= sqlc.narg('date_from')::timestamptz
    )
    AND (
      sqlc.narg('date_to')::timestamptz IS NULL
      OR m.created_at <= sqlc.narg('date_to')::timestamptz
    )
)
SELECT
  category_id,
  category_name,
  category_is_system,
  category_system_key,
  COUNT(*)::bigint AS movement_count,
  COUNT(*) FILTER (WHERE movement_type_key = 'income')::bigint AS income_count,
  COUNT(*) FILTER (WHERE movement_type_key = 'expense')::bigint AS expense_count,
  COALESCE(SUM(amount) FILTER (WHERE movement_type_key = 'income'), 0)::numeric AS income_total,
  COALESCE(SUM(amount) FILTER (WHERE movement_type_key = 'expense'), 0)::numeric AS expense_total,
  COALESCE(SUM(CASE
    WHEN movement_type_key = 'income' THEN amount
    WHEN movement_type_key = 'expense' THEN -amount
    ELSE 0
  END), 0)::numeric AS net_total
FROM filtered
GROUP BY category_id, category_name, category_is_system, category_system_key
ORDER BY expense_total DESC, movement_count DESC, category_name ASC
;

-- name: DeleteMovement :execrows
UPDATE movements
SET
  deleted_at = $2,
  updated_at = $3
WHERE id = $1
  AND deleted_at IS NULL
;
