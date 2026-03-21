package domain

import "github.com/google/uuid"

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