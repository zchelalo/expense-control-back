package login

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"github.com/zchelalo/expense-control-back/internal/shared/crypto/password"
	"github.com/zchelalo/expense-control-back/internal/shared/idgen"
	"github.com/zchelalo/expense-control-back/internal/shared/observability"
	"go.uber.org/zap"
)

type UseCase struct {
	store    ports.CredentialStore
	sessions ports.SessionRepository
	hasher   password.PasswordHasher
	issuer   ports.TokenIssuer
	ids      idgen.Generator
	clock    clock.Clock
	refreshTTL time.Duration
}

func New(
	store ports.CredentialStore,
	sessions ports.SessionRepository,
	hasher password.PasswordHasher,
	issuer ports.TokenIssuer,
	ids idgen.Generator,
	clock clock.Clock,
	refreshTTL time.Duration,
) *UseCase {
	return &UseCase{
		store:      store,
		sessions:   sessions,
		hasher:     hasher,
		issuer:     issuer,
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

	// Normalize and validate email
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		log.Warn("invalid email format",
			zap.String("stage", "validate_input"),
			zap.String("email", cmd.Email),
		)
		return Result{}, err
	}

	// Retrieve credential by email
	cred, err := uc.store.ByEmail(ctx, email)
	if err != nil {
		log.Warn("email not found",
			zap.String("stage", "retrieve_credential"),
			zap.String("email", email.String()),
		)
		return Result{}, ErrInvalidCredentials
	}

	log.Info("credential retrieved",
		zap.String("stage", "credential_retrieved"),
		zap.String("credential_id", cred.ID().String()),
	)

	// Compare password
	if err := uc.hasher.Compare(ctx, cmd.Password, cred.PasswordHash().String()); err != nil {
		log.Warn("invalid password",
			zap.String("stage", "compare_password"),
			zap.String("credential_id", cred.ID().String()),
		)
    return Result{}, ErrInvalidCredentials
	}

	log.Info("password verified",
		zap.String("stage", "password_verified"),
		zap.String("credential_id", cred.ID().String()),
	)

	// Get current time
	now := uc.clock.Now()

	// Create session ID and refresh token ID
	sessID, err := domain.NewSessionID(uc.ids.NewUUID())
	if err != nil {
		log.Warn("failed to create session ID",
			zap.String("stage", "create_session_id"),
			zap.Error(err),
		)
		return Result{}, err
	}
	refreshID, err := domain.NewRefreshTokenID(uc.ids.NewUUID())
	if err != nil {
		log.Warn("failed to create refresh token ID",
			zap.String("stage", "create_refresh_id"),
			zap.Error(err),
		)
		return Result{}, err
	}

	// Generate refresh token expiration time
	refreshExp := now.Add(uc.refreshTTL)

	// Create session
	sess, err := domain.NewSession(sessID, cred.ID(), refreshID, now, refreshExp)
	if err != nil {
		log.Warn("failed to create session",
			zap.String("stage", "create_session"),
			zap.Error(err),
		)
		return Result{}, err
	}

	// Store session
	if err := uc.sessions.Create(ctx, sess); err != nil {
		log.Warn("failed to store session",
			zap.String("stage", "store_session"),
			zap.Error(err),
		)
		return Result{}, err
	}

	log.Info("",
		zap.String("stage", "session_created"),
		zap.String("credential_id", cred.ID().String()),
		zap.String("session_id", sessID.String()),
	)

	// Issue tokens
	accessToken, accessTokenExp, err := uc.issuer.IssueAccess(ctx, cred.ID())
	if err != nil {
		log.Warn("failed to issue access token",
			zap.String("stage", "issue_access_token"),
			zap.Error(err),
		)
		return Result{}, err
	}
	refreshToken, refreshTokenExp, err := uc.issuer.IssueRefresh(ctx, sessID, cred.ID(), refreshID)
	if err != nil {
		log.Warn("failed to issue refresh token",
			zap.String("stage", "issue_refresh_token"),
			zap.Error(err),
		)
		return Result{}, err
	}

	accessFp := observability.TokenFingerprint(accessToken)
	refreshFp := observability.TokenFingerprint(refreshToken)

	log.Info("tokens issued successfully",
		zap.String("stage", "issue_tokens"),
		zap.String("credential_id", cred.ID().String()),
		zap.String("session_id", sessID.String()),
		zap.String("access_token_fp", accessFp),
		zap.String("refresh_token_fp", refreshFp),
	)

	return Result{
		SubjectID:     cred.ID().String(),
		AccessToken:   accessToken,
		AccessExpires: accessTokenExp,
		RefreshToken:  refreshToken,
		RefreshExpires: refreshTokenExp,
	}, nil
}