package domain

import "errors"

var (
	// VO validation / invariants
	ErrInvalidAmount              = errors.New("invalid amount")
	ErrInvalidMovementID          = errors.New("invalid movement id")
	ErrInvalidMovementTypeID      = errors.New("invalid movement type id")
	ErrInvalidCategoryID          = errors.New("invalid category id")
	ErrInvalidAccountID           = errors.New("invalid account id")
	ErrInvalidUserID              = errors.New("invalid user id")
	ErrInvalidDescription         = errors.New("invalid description")
	ErrInvalidMovementType        = errors.New("invalid movement type")
	ErrInvalidCategory            = errors.New("invalid category")
	ErrInvalidAccount             = errors.New("invalid account")
	ErrInsufficientAccountBalance = errors.New("insufficient account balance")
)
