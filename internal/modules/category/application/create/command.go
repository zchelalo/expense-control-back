package create

import "github.com/google/uuid"

type Command struct {
	UserID uuid.UUID
	Name   string
}
