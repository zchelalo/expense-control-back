package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
)

type Middleware struct {
  Log *zap.Logger
	AllowedOrigins string
	Verifier ports.TokenVerifier
}

func New(log *zap.Logger, allowedOrigins string, verifier ports.TokenVerifier) *Middleware {
  return &Middleware{Log: log, AllowedOrigins: allowedOrigins, Verifier: verifier}
}

func ApplyMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}