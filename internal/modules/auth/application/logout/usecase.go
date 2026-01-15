package logout

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
)

type UseCase struct {
	verifier ports.TokenVerifier
	sessions ports.SessionRepository
	clock    clock.Clock
}

func New(
	verifier ports.TokenVerifier,
	sessions ports.SessionRepository,
	clock clock.Clock,
) *UseCase {
	return &UseCase{
		verifier: verifier,
		sessions:   sessions,
		clock:      clock,
	}
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) error {
	if cmd.RefreshToken == "" {
		return ErrMissingRefreshToken
	}

	// Verify refresh token
	claims, err := uc.verifier.VerifyRefresh(ctx, cmd.RefreshToken)
	if err != nil {
		return ports.ErrTokenInvalid{Name: "refresh"}
	}

	// Check that the token belongs to the subject
	if claims.SubjectID != cmd.SubjectID {
		return ErrForbidden
	}

	now := uc.clock.Now()

	// Revoke session
	if err = uc.sessions.Revoke(ctx, claims.SessionID, now); err != nil {
		return err
	}

	return nil
}