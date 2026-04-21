-- name: UserExists :one
SELECT EXISTS (
  SELECT 1
  FROM users
  WHERE id = $1
    AND deleted_at IS NULL
)
;
