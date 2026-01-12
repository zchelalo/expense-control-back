package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
)

type TokenIssuer interface {
	IssueAccess(ctx context.Context, sub domain.SubjectID) (string, time.Time, error)
	IssueRefresh(
		ctx context.Context,
		sessionID domain.SessionID,
		sub domain.SubjectID,
		refreshID domain.RefreshTokenID,
	) (string, time.Time, error)
}