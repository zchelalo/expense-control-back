package create

import "github.com/google/uuid"

type Command struct {
	Amount         float64
	Description    string
	MovementTypeID uuid.UUID
	CategoryID     uuid.UUID
	AccountID      uuid.UUID
	UserID         uuid.UUID
}
