package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/refresh"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type refreshRequest struct {
	RefreshToken   string `json:"refresh_token"`
}

type refreshResponse struct {
	SubjectID     string `json:"subject_id"`
	AccessToken   string `json:"access_token"`
	AccessExpires string `json:"access_expires_at"`

	RefreshToken   string `json:"refresh_token,omitempty"`
	RefreshExpires string `json:"refresh_expires_at,omitempty"`
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	rid := middleware.RequestIDFrom(r.Context())

	var refreshToken string

	mode := strings.ToLower(strings.TrimSpace(r.Header.Get(middleware.HeaderAuthMode)))
	isMobile := mode == middleware.AuthModeMobile

	var req refreshRequest
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

	res, err := h.refreshUC.Execute(r.Context(), refresh.Command{
		RefreshToken: refreshToken,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := refreshResponse{
		SubjectID:     res.SubjectID,
		AccessToken:   res.AccessToken,
		AccessExpires: res.AccessExpires.UTC().Format(time.RFC3339),
	}

	if isMobile {
		resp.RefreshToken = res.RefreshToken
		resp.RefreshExpires = res.RefreshExpires.UTC().Format(time.RFC3339)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     CookieRefreshToken,
			Value:    res.RefreshToken,
			Path:     AuthCookiePath,
			HttpOnly: true,
			Secure:   h.secureCookies,
			SameSite: http.SameSiteLaxMode,
			Expires:  res.RefreshExpires,
			MaxAge:   int(time.Until(res.RefreshExpires).Seconds()),
		})
	}

	response.WriteJSON(w, http.StatusCreated, resp, nil, rid)
}