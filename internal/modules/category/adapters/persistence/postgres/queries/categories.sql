-- name: UpsertCategoryByName :one
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
  FALSE,
  NULL,
  sqlc.arg('created_at'),
  sqlc.arg('updated_at'),
  sqlc.arg('deleted_at')
)
ON CONFLICT (name) WHERE is_system = FALSE DO UPDATE
SET
  deleted_at = NULL
RETURNING id, name
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

-- name: GetUserCategoryByUserIDAndCategoryID :one
SELECT
  uc.user_id,
  uc.category_id,
  uc.created_at,
  uc.updated_at,
  uc.deleted_at,
  c.is_system
FROM user_categories uc
INNER JOIN categories c
  ON c.id = uc.category_id
 AND c.deleted_at IS NULL
WHERE uc.user_id = sqlc.arg('user_id')
  AND uc.category_id = sqlc.arg('category_id')
;

-- name: CategoryHasActiveMovementsForUser :one
SELECT EXISTS (
  SELECT 1
  FROM movements m
  WHERE m.user_id = sqlc.arg('user_id')
    AND m.category_id = sqlc.arg('category_id')
    AND m.deleted_at IS NULL
)
;

-- name: DeleteUserCategory :execrows
UPDATE user_categories uc
SET
  deleted_at = sqlc.arg('deleted_at'),
  updated_at = sqlc.arg('updated_at')
WHERE uc.user_id = sqlc.arg('user_id')
  AND uc.category_id = sqlc.arg('category_id')
  AND uc.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1
    FROM movements m
    WHERE m.user_id = sqlc.arg('user_id')
      AND m.category_id = sqlc.arg('category_id')
      AND m.deleted_at IS NULL
  )
;

-- name: ListCategoriesByUserIDAfter :many
SELECT
  uc.category_id AS id,
  c.name,
  uc.user_id,
  uc.created_at,
  uc.updated_at
FROM user_categories uc
INNER JOIN categories c
  ON c.id = uc.category_id
 AND c.deleted_at IS NULL
WHERE uc.user_id = sqlc.arg('user_id')
  AND uc.deleted_at IS NULL
  AND c.is_system = FALSE
  AND (
    sqlc.narg('cursor_created_at')::timestamptz IS NULL
    OR (uc.created_at, uc.category_id) < (
      sqlc.narg('cursor_created_at')::timestamptz,
      sqlc.narg('cursor_category_id')::uuid
    )
  )
  AND (sqlc.narg('name')::text IS NULL OR c.name ILIKE '%' || sqlc.narg('name') || '%')
ORDER BY uc.created_at DESC, uc.category_id DESC
LIMIT sqlc.arg('limit_count')
;

-- name: ListCategoriesByUserIDBefore :many
SELECT
  uc.category_id AS id,
  c.name,
  uc.user_id,
  uc.created_at,
  uc.updated_at
FROM user_categories uc
INNER JOIN categories c
  ON c.id = uc.category_id
 AND c.deleted_at IS NULL
WHERE uc.user_id = sqlc.arg('user_id')
  AND uc.deleted_at IS NULL
  AND c.is_system = FALSE
  AND (
    sqlc.narg('cursor_created_at')::timestamptz IS NULL
    OR (uc.created_at, uc.category_id) > (
      sqlc.narg('cursor_created_at')::timestamptz,
      sqlc.narg('cursor_category_id')::uuid
    )
  )
  AND (sqlc.narg('name')::text IS NULL OR c.name ILIKE '%' || sqlc.narg('name') || '%')
ORDER BY uc.created_at ASC, uc.category_id ASC
LIMIT sqlc.arg('limit_count')
;
