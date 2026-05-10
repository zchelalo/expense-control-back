package stats

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	UserID         uuid.UUID
	AccountID      *uuid.UUID
	CategoryID     *uuid.UUID
	MovementTypeID *uuid.UUID
	DateFrom       *time.Time
	DateTo         *time.Time
}
