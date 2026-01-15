package logout

import "errors"

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionAlreadyRevoked = errors.New("session already revoked")
var ErrMissingRefreshToken = errors.New("missing refresh token")
var ErrForbidden = errors.New("forbidden")