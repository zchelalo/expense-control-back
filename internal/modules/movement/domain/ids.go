package domain

import "github.com/google/uuid"

type MovementID struct{ value uuid.UUID }

func NewMovementID(v uuid.UUID) (MovementID, error) {
	if v == uuid.Nil {
		return MovementID{}, ErrInvalidMovementID
	}
	return MovementID{value: v}, nil
}

func (id MovementID) UUID() uuid.UUID { return id.value }
func (id MovementID) String() string  { return id.value.String() }

type MovementTypeID struct{ value uuid.UUID }

func NewMovementTypeID(v uuid.UUID) (MovementTypeID, error) {
	if v == uuid.Nil {
		return MovementTypeID{}, ErrInvalidMovementTypeID
	}
	return MovementTypeID{value: v}, nil
}

func (id MovementTypeID) UUID() uuid.UUID { return id.value }
func (id MovementTypeID) String() string  { return id.value.String() }

type CategoryID struct{ value uuid.UUID }

func NewCategoryID(v uuid.UUID) (CategoryID, error) {
	if v == uuid.Nil {
		return CategoryID{}, ErrInvalidCategoryID
	}
	return CategoryID{value: v}, nil
}

func (id CategoryID) UUID() uuid.UUID { return id.value }
func (id CategoryID) String() string  { return id.value.String() }

type AccountID struct{ value uuid.UUID }

func NewAccountID(v uuid.UUID) (AccountID, error) {
	if v == uuid.Nil {
		return AccountID{}, ErrInvalidAccountID
	}
	return AccountID{value: v}, nil
}

func (id AccountID) UUID() uuid.UUID { return id.value }
func (id AccountID) String() string  { return id.value.String() }

type UserID struct{ value uuid.UUID }

func NewUserID(v uuid.UUID) (UserID, error) {
	if v == uuid.Nil {
		return UserID{}, ErrInvalidUserID
	}
	return UserID{value: v}, nil
}

func (id UserID) UUID() uuid.UUID { return id.value }
func (id UserID) String() string  { return id.value.String() }
