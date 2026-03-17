package v1

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/logout"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type logoutRequest struct {
	RefreshToken   string `json:"refresh_token"`
}

type logoutResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	rid := middleware.RequestIDFrom(r.Context())

	var refreshToken string

	mode := strings.ToLower(strings.TrimSpace(r.Header.Get(middleware.HeaderAuthMode)))
	isMobile := mode == middleware.AuthModeMobile

	var req logoutRequest
	if isMobile {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.WriteError(w, http.StatusBadRequest, response.APIError{
				Code:    "invalid_json",
				Message: "invalid json body",
			}, rid)
			return
		}
	}

	if !isMobile {
		cookie, err := r.Cookie(CookieRefreshToken)
		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, response.APIError{
				Code:    "invalid_input",
				Message: "refresh token cookie is required",
			}, rid)
			return
		}
		refreshToken = cookie.Value
	} else {
		refreshToken = req.RefreshToken
	}

	refreshToken = middleware.StripBearer(refreshToken)
	if refreshToken == "" {
		response.WriteError(w, http.StatusUnauthorized, response.APIError{
			Code:    "invalid_input",
			Message: "refresh token is required",
		}, rid)
		return
	}

	subID, ok := middleware.SubjectIDFrom(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, response.APIError{
			Code:    "unauthorized",
			Message: "subject id not found in context",
		}, rid)
		return
	}

	err := h.logoutUC.Execute(r.Context(), logout.Command{
		SubjectID:    subID,
		RefreshToken: refreshToken,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := logoutResponse{
		Message: "successfully logged out",
	}

	if !isMobile {
		http.SetCookie(w, &http.Cookie{
			Name:     CookieAccessToken,
			Value:    "",
			Path:     AccessCookiePath,
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   h.secureCookies,
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     CookieRefreshToken,
			Value:    "",
			Path:     RefreshCookiePath,
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   h.secureCookies,
			SameSite: http.SameSiteLaxMode,
		})
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}