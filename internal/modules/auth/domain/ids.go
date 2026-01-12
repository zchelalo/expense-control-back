package domain

import "github.com/google/uuid"

type SubjectID struct{ value uuid.UUID }

func NewSubjectID(v uuid.UUID) (SubjectID, error) {
	if v == uuid.Nil {
		return SubjectID{}, ErrInvalidSubjectID
	}
	return SubjectID{value: v}, nil
}

func (id SubjectID) UUID() uuid.UUID { return id.value }
func (id SubjectID) String() string  { return id.value.String() }

type SessionID struct{ value uuid.UUID }

func NewSessionID(v uuid.UUID) (SessionID, error) {
	if v == uuid.Nil {
		return SessionID{}, ErrInvalidSessionID
	}
	return SessionID{value: v}, nil
}

func (id SessionID) UUID() uuid.UUID { return id.value }
func (id SessionID) String() string  { return id.value.String() }

type RefreshTokenID struct{ value uuid.UUID }

func NewRefreshTokenID(v uuid.UUID) (RefreshTokenID, error) {
	if v == uuid.Nil {
		return RefreshTokenID{}, ErrInvalidRefreshTokenID
	}
	return RefreshTokenID{value: v}, nil
}

func (id RefreshTokenID) UUID() uuid.UUID { return id.value }
func (id RefreshTokenID) String() string  { return id.value.String() }