package register

import (
	"context"
	"errors"
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

	// Check if email already exists
	_, err = uc.store.ByEmail(ctx, email)
	if err == nil {
		return Result{}, ErrEmailAlreadyExists
	}

	// Hash password and create password hash
	hash, err := uc.hasher.Hash(ctx, cmd.Password)
	if err != nil {
		return Result{}, err
	}
	passHash, err := domain.NewPasswordHashFromHash(hash)
	if err != nil {
		return Result{}, err
	}

	// Generate new subject ID
	sub, err := domain.NewSubjectID(uc.ids.NewUUID())
	if err != nil {
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
				return Result{}, ErrEmailAlreadyExists
			}
    }
    return Result{}, err
	}

	// Create session ID and refresh token ID
	sessID, _ := domain.NewSessionID(uc.ids.NewUUID())
	refreshID, _ := domain.NewRefreshTokenID(uc.ids.NewUUID())

	// Generate refresh token expiration time
	refreshExp := now.Add(uc.refreshTTL)

	// Create session
	sess, err := domain.NewSession(sessID, createdSub, refreshID, now, refreshExp)
	if err != nil {
		return Result{}, err
	}

	// Store session
	if err := uc.sessions.Create(ctx, sess); err != nil {
		return Result{}, err
	}

	// Issue tokens
	accessToken, accessExp, err := uc.issuer.IssueAccess(ctx, createdSub)
	if err != nil {
		return Result{}, err
	}
	refreshToken, refreshTokenExp, err := uc.issuer.IssueRefresh(ctx, sessID, createdSub, refreshID)
	if err != nil {
		return Result{}, err
	}

	return Result{
		SubjectID:      createdSub.String(),
		AccessToken:    accessToken,
		AccessExpires:  accessExp,
		RefreshToken:   refreshToken,
		RefreshExpires: refreshTokenExp,
	}, nil
}