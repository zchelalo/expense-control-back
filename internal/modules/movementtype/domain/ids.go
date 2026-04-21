package domain

import "github.com/google/uuid"

type MovementTypeID struct{ value uuid.UUID }

func NewMovementTypeID(v uuid.UUID) (MovementTypeID, error) {
	if v == uuid.Nil {
		return MovementTypeID{}, ErrInvalidMovementTypeID
	}

	return MovementTypeID{value: v}, nil
}

func (id MovementTypeID) UUID() uuid.UUID { return id.value }
func (id MovementTypeID) String() string  { return id.value.String() }
