package byid

import "github.com/google/uuid"

type Command struct {
	UserID    uuid.UUID
	AccountID uuid.UUID
}
