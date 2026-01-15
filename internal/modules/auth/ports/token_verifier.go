package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
)

type AccessClaims struct {
	SubjectID domain.SubjectID
	ExpiresAt time.Time
}

type RefreshClaims struct {
	SessionID domain.SessionID
	SubjectID domain.SubjectID
	RefreshID domain.RefreshTokenID
	ExpiresAt time.Time
}

type TokenVerifier interface {
	VerifyAccess(ctx context.Context, token string) (AccessClaims, error)
	VerifyRefresh(ctx context.Context, token string) (RefreshClaims, error)
}