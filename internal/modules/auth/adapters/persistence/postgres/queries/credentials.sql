-- name: CreateUser :one
INSERT INTO users (
  id,
  email,
  password,
  created_at,
  updated_at,
  deleted_at
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
  id
;

-- name: GetUserByEmail :one
SELECT
  id,
  email,
  password,
  created_at,
  updated_at,
  deleted_at
FROM users
WHERE email = $1
AND deleted_at IS NULL
;