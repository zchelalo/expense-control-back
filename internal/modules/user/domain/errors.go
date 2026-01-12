package domain

import "errors"

var (
	ErrInvalidUserID       = errors.New("invalid user id")
	ErrInvalidEmail        = errors.New("invalid email")
)