package v1

import (
	"errors"
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/login"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/logout"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/refresh"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/register"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/ports"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

func mapError(err error) (int, response.APIError) {
	// Registration errors
	if errors.Is(err, register.ErrEmailAlreadyExists) {
		return http.StatusConflict, response.APIError{
			Code:    "email_already_exists",
			Message: err.Error(),
		}
	}

	// Login errors
	if errors.Is(err, login.ErrInvalidCredentials) {
		return http.StatusUnauthorized, response.APIError{
			Code:    "invalid_credentials",
			Message: err.Error(),
		}
	}

	// Logout errors
	if errors.Is(err, logout.ErrSessionNotFound) {
		return http.StatusUnauthorized, response.APIError{
			Code:    "session_not_found",
			Message: err.Error(),
		}
	}
	if errors.Is(err, logout.ErrSessionAlreadyRevoked) {
		return http.StatusUnauthorized, response.APIError{
			Code:    "session_already_revoked",
			Message: err.Error(),
		}
	}
	if errors.Is(err, logout.ErrMissingRefreshToken) {
		return http.StatusUnauthorized, response.APIError{
			Code:    "missing_refresh_token",
			Message: err.Error(),
		}
	}
	if errors.Is(err, logout.ErrForbidden) {
		return http.StatusForbidden, response.APIError{
			Code:    "forbidden",
			Message: err.Error(),
		}
	}

	// Refresh errors
	if errors.Is(err, refresh.ErrSessionNotFound) {
		return http.StatusUnauthorized, response.APIError{
			Code:    "session_not_found",
			Message: err.Error(),
		}
	}
	if errors.Is(err, refresh.ErrSessionRevoked) {
		return http.StatusUnauthorized, response.APIError{
			Code:    "session_revoked",
			Message: err.Error(),
		}
	}
	if errors.Is(err, refresh.ErrMissingRefreshToken) {
		return http.StatusUnauthorized, response.APIError{
			Code:    "missing_refresh_token",
			Message: err.Error(),
		}
	}

	// Ports errors
	if errors.As(err, &ports.ErrTokenInvalid{}) {
		return http.StatusUnauthorized, response.APIError{
			Code:    "token_invalid",
			Message: err.Error(),
		}
	}

	// Domain validation
	if errors.Is(err, domain.ErrInvalidEmail) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_email",
			Message: err.Error(),
		}
	}

	// Default
	return http.StatusInternalServerError, response.APIError{
		Code:    "internal_error",
		Message: "internal server error",
	}
}