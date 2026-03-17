package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/register"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	SubjectID     string `json:"subject_id"`
	AccessToken   string `json:"access_token,omitempty"`
	AccessExpires string `json:"access_expires_at,omitempty"`

	RefreshToken   string `json:"refresh_token,omitempty"`
	RefreshExpires string `json:"refresh_expires_at,omitempty"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	rid := middleware.RequestIDFrom(r.Context())

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_json",
			Message: "invalid json body",
		}, rid)
		return
	}

	if strings.TrimSpace(req.Email) == "" || req.Password == "" {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_input",
			Message: "email and password are required",
		}, rid)
		return
	}

	res, err := h.registerUC.Execute(r.Context(), register.Command{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	mode := strings.ToLower(strings.TrimSpace(r.Header.Get(middleware.HeaderAuthMode)))
	isMobile := mode == middleware.AuthModeMobile

	resp := registerResponse{
		SubjectID: res.SubjectID,
	}

	if isMobile {
		resp.AccessToken = res.AccessToken
		resp.AccessExpires = res.AccessExpires.UTC().Format(time.RFC3339)
		resp.RefreshToken = res.RefreshToken
		resp.RefreshExpires = res.RefreshExpires.UTC().Format(time.RFC3339)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     CookieAccessToken,
			Value:    res.AccessToken,
			Path:     AccessCookiePath,
			HttpOnly: true,
			Secure:   h.secureCookies,
			SameSite: http.SameSiteLaxMode,
			Expires:  res.AccessExpires,
			MaxAge:   int(time.Until(res.AccessExpires).Seconds()),
		})

		http.SetCookie(w, &http.Cookie{
			Name:     CookieRefreshToken,
			Value:    res.RefreshToken,
			Path:     RefreshCookiePath,
			HttpOnly: true,
			Secure:   h.secureCookies,
			SameSite: http.SameSiteLaxMode,
			Expires:  res.RefreshExpires,
			MaxAge:   int(time.Until(res.RefreshExpires).Seconds()),
		})
	}

	response.WriteJSON(w, http.StatusCreated, resp, nil, rid)
}