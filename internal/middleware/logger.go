package middleware

import (
	"context"
	"net/http"

	"github.com/zchelalo/expense-control-back/pkg/observability"
	"go.uber.org/zap"
)

type ctxKeyLogger struct{}

func LoggerFrom(ctx context.Context) *zap.Logger {
	if v := ctx.Value(ctxKeyLogger{}); v != nil {
		if l, ok := v.(*zap.Logger); ok {
			return l
		}
	}
	return zap.NewNop()
}

func (m *Middleware) InjectLogger(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    base := m.Log
    l := observability.WithTrace(r.Context(), base)
    ctx := context.WithValue(r.Context(), ctxKeyLogger{}, l)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}