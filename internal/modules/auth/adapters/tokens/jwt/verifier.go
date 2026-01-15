package jwt

import (
	"context"
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
)

type Verifier struct {
	keys   *Keys
	issuer string
}

func NewVerifier(keys *Keys, issuer string) *Verifier {
	return &Verifier{keys: keys, issuer: issuer}
}

func (v *Verifier) VerifyAccess(_ context.Context, token string) (ports.AccessClaims, error) {
	parsed, err := jwtlib.ParseWithClaims(token, &AccessClaims{}, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodRSA); !ok {
			return nil, ports.ErrTokenInvalid{Name: "access"}
		}
		return v.keys.AccessPublic, nil
	},
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodRS256.Alg()}),
		jwtlib.WithIssuer(v.issuer),
	)
	if err != nil {
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return ports.AccessClaims{}, ports.ErrTokenExpired{Name: "access"}
		}
		return ports.AccessClaims{}, ports.ErrTokenInvalid{Name: "access"}
	}

	c, ok := parsed.Claims.(*AccessClaims)
	if !ok || !parsed.Valid {
		return ports.AccessClaims{}, ports.ErrTokenInvalid{Name: "access"}
	}

	parsedSub, err := uuid.Parse(c.Sub)
	if err != nil {
		return ports.AccessClaims{}, ports.ErrTokenInvalid{Name: "access"}
	}
	sub, err := domain.NewSubjectID(parsedSub)
	if err != nil {
		return ports.AccessClaims{}, ports.ErrTokenInvalid{Name: "access"}
	}

	var exp time.Time
	if c.ExpiresAt != nil {
		exp = c.ExpiresAt.Time
	}

	return ports.AccessClaims{
		SubjectID: sub,
		ExpiresAt: exp,
	}, nil
}

func (v *Verifier) VerifyRefresh(_ context.Context, token string) (ports.RefreshClaims, error) {
	parsed, err := jwtlib.ParseWithClaims(token, &RefreshClaims{}, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodRSA); !ok {
			return nil, ports.ErrTokenInvalid{Name: "refresh"}
		}
		return v.keys.RefreshPublic, nil
	}, jwtlib.WithValidMethods([]string{jwtlib.SigningMethodRS256.Alg()}), jwtlib.WithIssuer(v.issuer))
	if err != nil {
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return ports.RefreshClaims{}, ports.ErrTokenExpired{Name: "refresh"}
		}
		return ports.RefreshClaims{}, ports.ErrTokenInvalid{Name: "refresh"}
	}

	c, ok := parsed.Claims.(*RefreshClaims)
	if !ok || !parsed.Valid {
		return ports.RefreshClaims{}, ports.ErrTokenInvalid{Name: "refresh"}
	}

	parsedSid, err := uuid.Parse(c.SID)
	if err != nil {
		return ports.RefreshClaims{}, ports.ErrTokenInvalid{Name: "refresh"}
	}
	sid, err := domain.NewSessionID(parsedSid)
	if err != nil { return ports.RefreshClaims{}, ports.ErrTokenInvalid{Name: "refresh"} }

	parsedSub, err := uuid.Parse(c.Sub)
	if err != nil {
		return ports.RefreshClaims{}, ports.ErrTokenInvalid{Name: "refresh"}
	}
	sub, err := domain.NewSubjectID(parsedSub)
	if err != nil { return ports.RefreshClaims{}, ports.ErrTokenInvalid{Name: "refresh"} }

	parsedJti, err := uuid.Parse(c.JTI)
	if err != nil {
		return ports.RefreshClaims{}, ports.ErrTokenInvalid{Name: "refresh"}
	}
	jti, err := domain.NewRefreshTokenID(parsedJti)
	if err != nil { return ports.RefreshClaims{}, ports.ErrTokenInvalid{Name: "refresh"} }

	var exp time.Time
	if c.ExpiresAt != nil {
		exp = c.ExpiresAt.Time
	}

	return ports.RefreshClaims{
		SessionID: sid,
		SubjectID: sub,
		RefreshID: jti,
		ExpiresAt: exp,
	}, nil
}