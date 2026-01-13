-- name: CreateSession :one
INSERT INTO auth_sessions (
  id,
  user_id,
  refresh_jti,
  created_at,
  expires_at,
  revoked_at
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
  id,
  user_id,
  refresh_jti,
  created_at,
  expires_at,
  revoked_at
;

-- name: GetSessionByID :one
SELECT
  id,
  user_id,
  refresh_jti,
  created_at,
  expires_at,
  revoked_at
FROM auth_sessions
WHERE id = $1
AND revoked_at IS NULL
;

-- name: RotateSessionRefreshID :one
UPDATE auth_sessions
SET
  refresh_jti = $2,
  expires_at = $3
WHERE id = $1
AND revoked_at IS NULL
RETURNING
  id,
  user_id,
  refresh_jti,
  created_at,
  expires_at,
  revoked_at
;

-- name: RevokeSession :one
UPDATE auth_sessions
SET
  revoked_at = $2
WHERE id = $1
AND revoked_at IS NULL
RETURNING
  id,
  user_id,
  refresh_jti,
  created_at,
  expires_at,
  revoked_at
;