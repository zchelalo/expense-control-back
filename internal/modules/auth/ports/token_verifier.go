package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AccessClaims struct {
	SubjectID uuid.UUID
	ExpiresAt time.Time
}

type RefreshClaims struct {
	SessionID uuid.UUID
	SubjectID uuid.UUID
	RefreshID uuid.UUID
	ExpiresAt time.Time
}

type TokenVerifier interface {
	VerifyAccess(ctx context.Context, token string) (AccessClaims, error)
	VerifyRefresh(ctx context.Context, token string) (RefreshClaims, error)
}
