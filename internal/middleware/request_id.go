package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type ctxKeyRequestID struct{}

const RequestIDHeader = "X-Request-Id"

func (m *Middleware) RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(RequestIDHeader)
		if rid == "" {
			rid = newRequestID()
		}

		w.Header().Set(RequestIDHeader, rid)

		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFrom(ctx context.Context) string {
	if v := ctx.Value(ctxKeyRequestID{}); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func newRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}