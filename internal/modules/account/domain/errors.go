package domain

import "errors"

var (
	// VO validation / invariants
	ErrInvalidAccountID  = errors.New("invalid account id")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidName       = errors.New("invalid name")
	ErrInvalidBalance    = errors.New("invalid balance")
)