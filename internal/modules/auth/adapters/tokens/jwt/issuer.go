package jwt

import (
	"context"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
)

type Issuer struct {
	keys       *Keys
	clock      clock.Clock
	accessTTL  time.Duration
	refreshTTL time.Duration
	issuer     string
}

func NewIssuer(keys *Keys, clock clock.Clock, issuer string, accessTTL, refreshTTL time.Duration) *Issuer {
	return &Issuer{
		keys:       keys,
		clock:      clock,
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (i *Issuer) IssueAccess(_ context.Context, sub uuid.UUID) (string, time.Time, error) {
	now := i.clock.Now()
	exp := now.Add(i.accessTTL)

	claims := AccessClaims{
		Sub: sub.String(),
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    i.issuer,
			Subject:   sub.String(),
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(exp),
		},
	}

	tok := jwtlib.NewWithClaims(jwtlib.SigningMethodRS256, claims)
	s, err := tok.SignedString(i.keys.AccessPrivate)
	if err != nil {
		return "", time.Time{}, err
	}
	return s, exp, nil
}

func (i *Issuer) IssueRefresh(_ context.Context, sessionID uuid.UUID, sub uuid.UUID, refreshID uuid.UUID) (string, time.Time, error) {
	now := i.clock.Now()
	exp := now.Add(i.refreshTTL)

	claims := RefreshClaims{
		SID: sessionID.String(),
		Sub: sub.String(),
		JTI: refreshID.String(),
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    i.issuer,
			Subject:   sub.String(),
			ID:        refreshID.String(),
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(exp),
		},
	}

	tok := jwtlib.NewWithClaims(jwtlib.SigningMethodRS256, claims)
	s, err := tok.SignedString(i.keys.RefreshPrivate)
	if err != nil {
		return "", time.Time{}, err
	}
	return s, exp, nil
}