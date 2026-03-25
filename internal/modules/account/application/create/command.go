package create

import "github.com/google/uuid"

type Command struct {
	Name    string
	Balance float64
	UserID  uuid.UUID
}
