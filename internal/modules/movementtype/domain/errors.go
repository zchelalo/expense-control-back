package domain

import "errors"

var (
	ErrInvalidMovementTypeID = errors.New("invalid movement type id")
	ErrInvalidMovementType   = errors.New("invalid movement type")
)
