package list

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	UserID     uuid.UUID
	CreatedAt  *time.Time
	CategoryID *uuid.UUID
	Limit      int
	IsBefore   bool
}
