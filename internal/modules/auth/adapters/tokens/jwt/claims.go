package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	SID string `json:"sid"`
	Sub string `json:"sub"`
	JTI string `json:"jti"`
	jwt.RegisteredClaims
}
