package register

import (
	"context"
	"errors"
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

	// Check if email already exists
	_, err = uc.store.ByEmail(ctx, email)
	if err == nil {
		log.Warn("email already exists",
			zap.String("stage", "check_existing"),
			zap.String("email", email.String()),
		)
		return Result{}, ErrEmailAlreadyExists
	}

	// Hash password and create password hash
	hash, err := uc.hasher.Hash(ctx, cmd.Password)
	if err != nil {
		log.Warn("failed to hash password",
			zap.String("stage", "hash_password"),
			zap.String("email", email.String()),
			zap.Error(err),
		)
		return Result{}, err
	}
	passHash, err := domain.NewPasswordHashFromHash(hash)
	if err != nil {
		log.Warn("invalid password hash generated",
			zap.String("stage", "create_pass_hash"),
			zap.String("email", email.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	// Generate new subject ID
	sub, err := domain.NewSubjectID(uc.ids.NewUUID())
	if err != nil {
		log.Warn("failed to generate subject ID",
			zap.String("stage", "generate_subject_id"),
			zap.String("email", email.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	// Get current time
	now := uc.clock.Now()

	// Create account
	account := domain.NewAccount(sub, email, passHash, now)

	// Store account
	createdSub, err := uc.store.CreateAccount(ctx, account)
	if err != nil {
		// Handle already exists error
    var exists ports.ErrAlreadyExists
    if errors.As(err, &exists) {
			switch exists.Name {
			case "email":
				log.Warn("email already exists on create",
					zap.String("stage", "store_account"),
					zap.String("email", email.String()),
				)
				return Result{}, ErrEmailAlreadyExists
			}
    }

		log.Warn("failed to store account",
			zap.String("stage", "store_account"),
			zap.String("email", email.String()),
			zap.Error(err),
		)
    return Result{}, err
	}

	log.Info("account registered",
		zap.String("stage", "account_created"),
		zap.String("subject_id", createdSub.String()),
		zap.String("email", email.String()),
	)

	// Create session ID and refresh token ID
	sessID, err := domain.NewSessionID(uc.ids.NewUUID())
	if err != nil {
		log.Warn("failed to generate session ID",
			zap.String("stage", "generate_session_id"),
			zap.String("subject_id", createdSub.String()),
			zap.Error(err),
		)
		return Result{}, err
	}
	refreshID, err := domain.NewRefreshTokenID(uc.ids.NewUUID())
	if err != nil {
		log.Warn("failed to generate refresh token ID",
			zap.String("stage", "generate_refresh_id"),
			zap.String("subject_id", createdSub.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	// Generate refresh token expiration time
	refreshExp := now.Add(uc.refreshTTL)

	// Create session
	sess, err := domain.NewSession(sessID, createdSub, refreshID, now, refreshExp)
	if err != nil {
		log.Warn("failed to create session",
			zap.String("stage", "create_session"),
			zap.String("subject_id", createdSub.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	// Store session
	if err := uc.sessions.Create(ctx, sess); err != nil {
		log.Warn("failed to store session",
			zap.String("stage", "store_session"),
			zap.String("subject_id", createdSub.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	log.Info("session stored successfully",
		zap.String("stage", "store_session"),
		zap.String("subject_id", createdSub.String()),
		zap.String("session_id", sessID.String()),
	)

	// Issue tokens
	accessToken, accessTokenExp, err := uc.issuer.IssueAccess(ctx, createdSub)
	if err != nil {
		log.Warn("issue access token failed",
			zap.String("stage", "issue_access"),
			zap.String("subject_id", createdSub.String()),
			zap.String("session_id", sessID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}
	refreshToken, refreshTokenExp, err := uc.issuer.IssueRefresh(ctx, sessID, createdSub, refreshID)
	if err != nil {
		log.Warn("issue refresh token failed",
			zap.String("stage", "issue_refresh"),
			zap.String("subject_id", createdSub.String()),
			zap.String("session_id", sessID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	accessFp := observability.TokenFingerprint(accessToken)
	refreshFp := observability.TokenFingerprint(refreshToken)

	log.Info("tokens issued successfully",
		zap.String("stage", "issue_tokens"),
		zap.String("subject_id", createdSub.String()),
		zap.String("session_id", sessID.String()),
		zap.String("access_token_fp", accessFp),
		zap.String("refresh_token_fp", refreshFp),
	)

	return Result{
		SubjectID:      createdSub.String(),
		AccessToken:    accessToken,
		AccessExpires:  accessTokenExp,
		RefreshToken:   refreshToken,
		RefreshExpires: refreshTokenExp,
	}, nil
}