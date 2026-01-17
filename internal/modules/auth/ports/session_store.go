package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, s domain.Session) error
	ByID(ctx context.Context, id domain.SessionID) (domain.Session, error)
	ValidateAndRotateRefresh(ctx context.Context,
    sessionID domain.SessionID,
    expectedCurrent domain.RefreshTokenID,
    newRefreshID domain.RefreshTokenID,
    newExp time.Time,
	) (ok bool, err error)
	Revoke(ctx context.Context, sessionID domain.SessionID, now time.Time) error
}