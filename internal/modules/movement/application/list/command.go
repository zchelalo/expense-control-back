package list

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	UserID         uuid.UUID
	AccountID      *uuid.UUID
	CategoryID     *uuid.UUID
	MovementTypeID *uuid.UUID
	CreatedAt      *time.Time
	MovementID     *uuid.UUID
	Limit          int
	IsBefore       bool
}
