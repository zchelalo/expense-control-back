package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
)

type RefreshClaims struct {
	SessionID domain.SessionID
	SubjectID domain.SubjectID
	RefreshID domain.RefreshTokenID
	ExpiresAt time.Time
}

type TokenVerifier interface {
	VerifyRefresh(ctx context.Context, token string) (RefreshClaims, error)
}