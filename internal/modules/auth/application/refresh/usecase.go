package refresh

import (
	"context"
	"errors"
	"time"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"github.com/zchelalo/expense-control-back/internal/shared/idgen"
	"github.com/zchelalo/expense-control-back/internal/shared/observability"
	"go.uber.org/zap"
)

type UseCase struct {
	verifier ports.TokenVerifier
	issuer   ports.TokenIssuer
	sessions ports.SessionRepository
	ids      idgen.Generator
	clock    clock.Clock
	refreshTTL time.Duration
}

func New(
	verifier ports.TokenVerifier,
	issuer ports.TokenIssuer,
	sessions ports.SessionRepository,
	ids idgen.Generator,
	clock clock.Clock,
	refreshTTL time.Duration,
) *UseCase {
	return &UseCase{
		verifier: verifier,
		issuer:     issuer,
		sessions:   sessions,
		ids:        ids,
		clock:      clock,
		refreshTTL: refreshTTL,
	}
}

type Result struct {
	SubjectID     string
	AccessToken   string
	AccessExpires time.Time
	RefreshToken  string
	RefreshExpires time.Time
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	if cmd.RefreshToken == "" {
		log.Warn("missing refresh token",
			zap.String("stage", "validate_input"),
		)
		return Result{}, ErrMissingRefreshToken
	}

	fp := observability.TokenFingerprint(cmd.RefreshToken)

	// Verify refresh token
	claims, err := uc.verifier.VerifyRefresh(ctx, cmd.RefreshToken)
	if err != nil {
    kind := classifyTokenError(err)

		log.Warn("refresh verification failed",
			zap.String("stage", "verify_refresh"),
			zap.String("kind", kind),
			zap.String("token_fp", fp),
			zap.Error(err),
		)

		return Result{}, ports.ErrTokenInvalid{Name: "refresh"}
	}

	log.Info("refresh verified",
		zap.String("stage", "verify_refresh_ok"),
		zap.String("subject_id", claims.SubjectID.String()),
		zap.String("session_id", claims.SessionID.String()),
		zap.String("token_fp", fp),
	)

	now := uc.clock.Now()

	// Issue tokens
	newRefreshToken, newRefreshTokenExp, err := uc.rotateRefreshToken(ctx, fp, claims.RefreshID, claims.SessionID, claims.SubjectID, now)
	if err != nil {
		if errors.Is(err, ports.ErrSessionRefreshMismatch) {
			return Result{}, ports.ErrTokenInvalid{Name: "refresh"}
    }

    var tokenInvalid ports.ErrTokenInvalid
    if errors.As(err, &tokenInvalid) {
			return Result{}, tokenInvalid
    }

    log.Error("refresh rotation failed",
			zap.String("stage", "rotate_refresh"),
			zap.String("subject_id", claims.SubjectID.String()),
			zap.String("session_id", claims.SessionID.String()),
			zap.Error(err),
    )
    return Result{}, err
	}

	accessToken, accessTokenExp, err := uc.issuer.IssueAccess(ctx, claims.SubjectID)
	if err != nil {
		log.Error("issue access token failed",
			zap.String("stage", "issue_access"),
			zap.String("subject_id", claims.SubjectID.String()),
			zap.String("session_id", claims.SessionID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	res := Result{
		SubjectID:      claims.SubjectID.String(),
		AccessToken:    accessToken,
		AccessExpires:  accessTokenExp,
		RefreshToken:   newRefreshToken,
		RefreshExpires: newRefreshTokenExp,
	}

	log.Info("refresh finished",
		zap.String("stage", "done"),
		zap.String("subject_id", claims.SubjectID.String()),
		zap.String("session_id", claims.SessionID.String()),
		zap.String("token_fp", fp),
		zap.Bool("rotated", true),
		zap.Time("access_exp", accessTokenExp),
		zap.Time("new_refresh_exp", newRefreshTokenExp),
		zap.Time("now", now),
	)

	return res, nil
}

func (uc *UseCase) rotateRefreshToken(
	ctx context.Context,
	tokenFP string,
	expectedCurrent domain.RefreshTokenID,
	sessionID domain.SessionID,
	subjectID domain.SubjectID,
	now time.Time,
) (string, time.Time, error) {
	log := middleware.LoggerFrom(ctx)

	newUUID := uc.ids.NewUUID()
	newRefreshID, err := domain.NewRefreshTokenID(newUUID)
	if err != nil {
		log.Error("generate refresh id failed",
			zap.String("stage", "generate_refresh_id"),
			zap.String("session_id", sessionID.String()),
			zap.String("subject_id", subjectID.String()),
			zap.Error(err),
		)
		return "", time.Time{}, err
	}

	newExp := uc.clock.Now().Add(uc.refreshTTL)

	ok, err := uc.sessions.ValidateAndRotateRefresh(ctx, sessionID, expectedCurrent, newRefreshID, newExp)
	if err != nil {
		if errors.Is(err, ports.ErrSessionRefreshMismatch) {
			// Mismatch detected
			log.Warn("refresh reuse detected (mismatch)",
				zap.String("stage", "reuse_detection"),
				zap.Bool("security_event", true),
				zap.String("session_id", sessionID.String()),
				zap.String("subject_id", subjectID.String()),
				zap.String("token_fp", tokenFP),
			)

			if rerr := uc.sessions.Revoke(ctx, sessionID, now); rerr != nil {
				log.Error("failed to revoke session after refresh mismatch",
					zap.String("stage", "revoke_session"),
					zap.String("session_id", sessionID.String()),
					zap.String("subject_id", subjectID.String()),
					zap.Error(rerr),
				)
			}

			return "", time.Time{}, ports.ErrSessionRefreshMismatch
		}

		log.Error("persist refresh rotation failed",
			zap.String("stage", "persist_rotation"),
			zap.String("session_id", sessionID.String()),
			zap.String("subject_id", subjectID.String()),
			zap.Time("new_exp", newExp),
			zap.Error(err),
    )
    return "", time.Time{}, err
	}

	if !ok {
		// Session invalid (revoked or expired)
		log.Info("refresh rotation rejected (session invalid)",
			zap.String("stage", "session_invalid"),
			zap.String("session_id", sessionID.String()),
			zap.String("subject_id", subjectID.String()),
			zap.String("token_fp", tokenFP),
    )

		return "", time.Time{}, ports.ErrTokenInvalid{Name: "refresh"}
	}

	refreshToken, refreshTokenExp, err := uc.issuer.IssueRefresh(ctx, sessionID, subjectID, newRefreshID)
	if err != nil {
		log.Error("issue refresh token failed",
			zap.String("stage", "issue_refresh"),
			zap.String("session_id", sessionID.String()),
			zap.String("subject_id", subjectID.String()),
			zap.Time("new_exp", newExp),
			zap.Error(err),
		)
		return "", time.Time{}, err
	}

	log.Info("refresh rotated",
		zap.String("stage", "rotate_refresh_ok"),
		zap.String("session_id", sessionID.String()),
		zap.String("subject_id", subjectID.String()),
		zap.Time("new_refresh_exp", refreshTokenExp),
	)

	return refreshToken, refreshTokenExp, nil
}

func classifyTokenError(err error) string {
	var expired ports.ErrTokenExpired
	var sig ports.ErrTokenSignatureInvalid
	var malformed ports.ErrTokenMalformed

	switch {
	case errors.As(err, &expired):
		return "expired"
	case errors.As(err, &sig):
		return "bad_signature"
	case errors.As(err, &malformed):
		return "malformed"
	default:
		return "invalid"
	}
}