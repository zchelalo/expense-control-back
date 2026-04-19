package jwt

import (
	"context"
	"errors"
	"strings"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
	uuidparse "github.com/zchelalo/expense-control-back/pkg/parse"
)

type Verifier struct {
	keys   *Keys
	issuer string
}

func NewVerifier(keys *Keys, issuer string) *Verifier {
	return &Verifier{keys: keys, issuer: issuer}
}

func (v *Verifier) VerifyAccess(_ context.Context, token string) (ports.AccessClaims, error) {
	parsed, err := jwtlib.ParseWithClaims(
		token,
		&AccessClaims{},
		func(t *jwtlib.Token) (any, error) {
			if _, ok := t.Method.(*jwtlib.SigningMethodRSA); !ok {
				return nil, ports.ErrTokenSignatureInvalid{Name: "access"}
			}
			return v.keys.AccessPublic, nil
		},
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodRS256.Alg()}),
		jwtlib.WithIssuer(v.issuer),
	)
	if err != nil {
		return ports.AccessClaims{}, mapJWTParseError(err, "access")
	}

	c, ok := parsed.Claims.(*AccessClaims)
	if !ok || !parsed.Valid {
		return ports.AccessClaims{}, ports.ErrTokenInvalid{Name: "access"}
	}

	parsedSub, err := uuidparse.UUID(c.Sub)
	if err != nil {
		return ports.AccessClaims{}, ports.ErrTokenMalformed{Name: "access"}
	}

	var exp time.Time
	if c.ExpiresAt != nil {
		exp = c.ExpiresAt.Time
	}

	return ports.AccessClaims{
		SubjectID: parsedSub,
		ExpiresAt: exp,
	}, nil
}

func (v *Verifier) VerifyRefresh(_ context.Context, token string) (ports.RefreshClaims, error) {
	parsed, err := jwtlib.ParseWithClaims(
		token,
		&RefreshClaims{},
		func(t *jwtlib.Token) (any, error) {
			if _, ok := t.Method.(*jwtlib.SigningMethodRSA); !ok {
				return nil, ports.ErrTokenSignatureInvalid{Name: "refresh"}
			}
			return v.keys.RefreshPublic, nil
		},
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodRS256.Alg()}),
		jwtlib.WithIssuer(v.issuer),
	)
	if err != nil {
		return ports.RefreshClaims{}, mapJWTParseError(err, "refresh")
	}

	c, ok := parsed.Claims.(*RefreshClaims)
	if !ok || !parsed.Valid {
		return ports.RefreshClaims{}, ports.ErrTokenInvalid{Name: "refresh"}
	}

	parsedSid, err := uuidparse.UUID(c.SID)
	if err != nil {
		return ports.RefreshClaims{}, ports.ErrTokenMalformed{Name: "refresh"}
	}

	parsedSub, err := uuidparse.UUID(c.Sub)
	if err != nil {
		return ports.RefreshClaims{}, ports.ErrTokenMalformed{Name: "refresh"}
	}

	parsedJti, err := uuidparse.UUID(c.JTI)
	if err != nil {
		return ports.RefreshClaims{}, ports.ErrTokenMalformed{Name: "refresh"}
	}

	var exp time.Time
	if c.ExpiresAt != nil {
		exp = c.ExpiresAt.Time
	}

	return ports.RefreshClaims{
		SessionID: parsedSid,
		SubjectID: parsedSub,
		RefreshID: parsedJti,
		ExpiresAt: exp,
	}, nil
}

func mapJWTParseError(err error, name string) error {
	if errors.Is(err, jwtlib.ErrTokenExpired) {
		return ports.ErrTokenExpired{Name: name}
	}

	if errors.Is(err, jwtlib.ErrTokenNotValidYet) {
		return ports.ErrTokenInvalid{Name: name}
	}

	var sig ports.ErrTokenSignatureInvalid
	if errors.As(err, &sig) {
		return sig
	}
	var malformed ports.ErrTokenMalformed
	if errors.As(err, &malformed) {
		return malformed
	}
	var expired ports.ErrTokenExpired
	if errors.As(err, &expired) {
		return expired
	}

	msg := strings.ToLower(err.Error())

	if strings.Contains(msg, "invalid number of segments") ||
		strings.Contains(msg, "token contains an invalid number of segments") ||
		strings.Contains(msg, "illegal base64") ||
		strings.Contains(msg, "invalid character") ||
		strings.Contains(msg, "cannot parse") {
		return ports.ErrTokenMalformed{Name: name}
	}

	if strings.Contains(msg, "signature is invalid") ||
		strings.Contains(msg, "verification error") ||
		strings.Contains(msg, "crypto/rsa") ||
		strings.Contains(msg, "key is invalid") {
		return ports.ErrTokenSignatureInvalid{Name: name}
	}

	return ports.ErrTokenInvalid{Name: name}
}
