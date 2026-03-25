package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenIssuer interface {
	IssueAccess(ctx context.Context, sub uuid.UUID) (string, time.Time, error)
	IssueRefresh(
		ctx context.Context,
		sessionID uuid.UUID,
		sub uuid.UUID,
		refreshID uuid.UUID,
	) (string, time.Time, error)
}
