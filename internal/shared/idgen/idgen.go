package idgen

import "github.com/google/uuid"

type Generator interface {
	NewUUID() uuid.UUID
}