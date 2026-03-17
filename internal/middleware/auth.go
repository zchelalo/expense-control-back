package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
	"github.com/zchelalo/expense-control-back/pkg/response"
	"go.uber.org/zap"
)

type ctxSubjectID struct{}
var subjectIDKey = ctxSubjectID{}

const HeaderAuthMode = "X-Auth-Mode"
const AuthModeMobile = "mobile"
const HeaderAccessToken = "X-Access-Token"
const CookieAccessToken = "access_token"

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := RequestIDFrom(r.Context())

		mode := strings.ToLower(strings.TrimSpace(r.Header.Get(HeaderAuthMode)))
		isMobile := mode == AuthModeMobile

		token := extractAccessToken(r, isMobile)
		if token == "" {
			m.Log.Warn("missing access token")
			response.WriteError(w, http.StatusUnauthorized, response.APIError{
				Code:    "missing_access_token",
				Message: "missing access token",
			}, rid)
			return
		}

		claims, err := m.Verifier.VerifyAccess(r.Context(), token)
		if err != nil {
			m.Log.Warn("invalid access token", zap.Error(err))
			response.WriteError(w, http.StatusUnauthorized, response.APIError{
				Code:    "invalid_access_token",
				Message: "invalid access token",
			}, rid)
			return
		}

		ctx := context.WithValue(r.Context(), subjectIDKey, claims.SubjectID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractAccessToken(r *http.Request, isMobile bool) string {
	if isMobile {
		if t := strings.TrimSpace(r.Header.Get(HeaderAccessToken)); t != "" {
			return StripBearer(t)
		}
	}

	if !isMobile {
		if cookie, err := r.Cookie(CookieAccessToken); err == nil {
			return strings.TrimSpace(cookie.Value)
		}
	}

	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return ""
	}

	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}

	if isMobile {
		return auth
	}

	return ""
}

func StripBearer(v string) string {
	v = strings.TrimSpace(v)
	if strings.HasPrefix(strings.ToLower(v), "bearer ") {
		return strings.TrimSpace(v[7:])
	}
	return v
}

func SubjectIDFrom(ctx context.Context) (domain.SubjectID, bool) {
	v := ctx.Value(subjectIDKey)
	if v == nil {
		return domain.SubjectID{}, false
	}
	sub, ok := v.(domain.SubjectID)
	return sub, ok
}