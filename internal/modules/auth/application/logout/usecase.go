package logout

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"go.uber.org/zap"
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
	log := middleware.LoggerFrom(ctx)

	if cmd.RefreshToken == "" {
		log.Warn("missing refresh token in logout request",
			zap.String("stage", "validate_input"),
		)
		return ErrMissingRefreshToken
	}

	// Verify refresh token
	claims, err := uc.verifier.VerifyRefresh(ctx, cmd.RefreshToken)
	if err != nil {
		log.Warn("invalid refresh token",
			zap.String("stage", "verify_refresh_token"),
			zap.Error(err),
		)
		return ports.ErrTokenInvalid{Name: "refresh"}
	}

	log.Info("refresh token verified successfully",
		zap.String("stage", "verify_refresh_token"),
		zap.String("subject_id", claims.SubjectID.String()),
		zap.String("session_id", claims.SessionID.String()),
	)

	// Check that the token belongs to the subject
	if claims.SubjectID != cmd.SubjectID {
		log.Warn("refresh token subject does not match command subject",
			zap.String("stage", "validate_token_subject"),
			zap.String("token_subject_id", claims.SubjectID.String()),
			zap.String("command_subject_id", cmd.SubjectID.String()),
		)
		return ErrForbidden
	}

	now := uc.clock.Now()

	// Revoke session
	if err = uc.sessions.Revoke(ctx, claims.SessionID, now); err != nil {
		log.Error("failed to revoke session during logout",
			zap.String("stage", "revoke_session"),
			zap.String("subject_id", claims.SubjectID.String()),
			zap.String("session_id", claims.SessionID.String()),
			zap.Error(err),
		)
		return err
	}

	log.Info("session revoked successfully",
		zap.String("stage", "revoke_session"),
		zap.String("subject_id", claims.SubjectID.String()),
		zap.String("session_id", claims.SessionID.String()),
	)

	return nil
}