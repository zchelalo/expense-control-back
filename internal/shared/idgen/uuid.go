package idgen

import "github.com/google/uuid"


type UUIDGenerator struct{}

func NewGenerator() *UUIDGenerator { return &UUIDGenerator{} }

func (g *UUIDGenerator) NewUUID() uuid.UUID {
	return uuid.New()
}