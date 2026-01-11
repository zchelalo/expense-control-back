package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

type Middleware struct {
  Log *zap.Logger
	AllowedOrigins string
}

func New(log *zap.Logger, allowedOrigins string) *Middleware {
  return &Middleware{Log: log, AllowedOrigins: allowedOrigins}
}

func ApplyMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}