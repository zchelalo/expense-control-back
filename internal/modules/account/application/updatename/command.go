package updatename

import "github.com/google/uuid"

type Command struct {
	UserID    uuid.UUID
	AccountID uuid.UUID
	Name      string
}
