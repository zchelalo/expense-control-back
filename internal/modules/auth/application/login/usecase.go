package login

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"github.com/zchelalo/expense-control-back/internal/shared/crypto/password"
	"github.com/zchelalo/expense-control-back/internal/shared/idgen"
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
	// Normalize and validate email
	email, err := domain.NewEmail(cmd.Email)
	if err != nil {
		return Result{}, err
	}

	// Retrieve credential by email
	cred, err := uc.store.ByEmail(ctx, email)
	if err != nil {
		return Result{}, ErrInvalidCredentials
	}

	// Compare password
	if err := uc.hasher.Compare(ctx, cmd.Password, cred.PasswordHash().String()); err != nil {
    return Result{}, ErrInvalidCredentials
	}

	// Get current time
	now := uc.clock.Now()

	// Create session ID and refresh token ID
	sessID, _ := domain.NewSessionID(uc.ids.NewUUID())
	refreshID, _ := domain.NewRefreshTokenID(uc.ids.NewUUID())

	// Generate refresh token expiration time
	refreshExp := now.Add(uc.refreshTTL)

	// Create session
	sess, err := domain.NewSession(sessID, cred.ID(), refreshID, now, refreshExp)
	if err != nil {
		return Result{}, err
	}

	// Store session
	if err := uc.sessions.Create(ctx, sess); err != nil {
		return Result{}, err
	}

	// Issue tokens
	accessToken, accessTokenExp, err := uc.issuer.IssueAccess(ctx, cred.ID())
	if err != nil {
		return Result{}, err
	}
	refreshToken, refreshTokenExp, err := uc.issuer.IssueRefresh(ctx, sessID, cred.ID(), refreshID)
	if err != nil {
		return Result{}, err
	}

	return Result{
		SubjectID:     cred.ID().String(),
		AccessToken:   accessToken,
		AccessExpires: accessTokenExp,
		RefreshToken:  refreshToken,
		RefreshExpires: refreshTokenExp,
	}, nil
}