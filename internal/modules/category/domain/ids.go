package domain

import "github.com/google/uuid"

type CategoryID struct{ value uuid.UUID }

func NewCategoryID(v uuid.UUID) (CategoryID, error) {
	if v == uuid.Nil {
		return CategoryID{}, ErrInvalidCategoryID
	}

	return CategoryID{value: v}, nil
}

func (id CategoryID) UUID() uuid.UUID { return id.value }
func (id CategoryID) String() string  { return id.value.String() }

type UserID struct{ value uuid.UUID }

func NewUserID(v uuid.UUID) (UserID, error) {
	if v == uuid.Nil {
		return UserID{}, ErrInvalidUserID
	}

	return UserID{value: v}, nil
}

func (id UserID) UUID() uuid.UUID { return id.value }
func (id UserID) String() string  { return id.value.String() }
