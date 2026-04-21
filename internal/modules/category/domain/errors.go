package domain

import "errors"

var (
	ErrInvalidCategoryID = errors.New("invalid category id")
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidName       = errors.New("invalid name")
	ErrInvalidCategory   = errors.New("invalid category")
)
