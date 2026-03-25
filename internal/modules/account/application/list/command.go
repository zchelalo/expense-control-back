package list

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	UserID    uuid.UUID
	Name      *string
	CreatedAt *time.Time
	AccountID *uuid.UUID
	Limit     int
	IsBefore  bool
}
