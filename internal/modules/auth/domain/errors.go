package domain

import "errors"

var (
	// VO validation / invariants
	ErrInvalidSubjectID     = errors.New("invalid subject id")
	ErrInvalidSessionID     = errors.New("invalid session id")
	ErrInvalidRefreshTokenID = errors.New("invalid refresh token id")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidPasswordHash  = errors.New("invalid password hash")

	// Domain state errors (session lifecycle)
	ErrSessionRevoked = errors.New("session is revoked")
	ErrSessionExpired = errors.New("session is expired")
)