package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, s domain.Session) error
	ByID(ctx context.Context, id domain.SessionID) (domain.Session, error)
	RotateRefresh(ctx context.Context, sessionID domain.SessionID, newRefreshID domain.RefreshTokenID, newExp time.Time) error
	Revoke(ctx context.Context, sessionID domain.SessionID, now time.Time) error
}