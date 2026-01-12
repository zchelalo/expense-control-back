package v1

import (
	"errors"
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/register"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
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