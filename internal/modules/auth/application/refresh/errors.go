package refresh

import "errors"

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionRevoked = errors.New("session is revoked")
var ErrMissingRefreshToken = errors.New("missing refresh token")
