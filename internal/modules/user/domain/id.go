package domain

import "github.com/google/uuid"

type UserID struct{ value uuid.UUID }

func NewUserID(v uuid.UUID) (UserID, error) {
	if v == uuid.Nil {
		return UserID{}, ErrInvalidUserID
	}
	return UserID{value: v}, nil
}

func (id UserID) UUID() uuid.UUID { return id.value }
func (id UserID) String() string  { return id.value.String() }