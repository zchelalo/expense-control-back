package v1

import (
	"errors"
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/modules/account/application/byid"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/list"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/account/ports"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

func mapError(err error) (int, response.APIError) {
	// Ports errors
	if errors.As(err, &ports.ErrNotFound{}) {
		return http.StatusNotFound, response.APIError{
			Code:    "not_found",
			Message: err.Error(),
		}
	}

	// Domain validation
	if errors.Is(err, domain.ErrInvalidAccountID) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_account_id",
			Message: err.Error(),
		}
	}
	if errors.Is(err, domain.ErrInvalidUserID) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_user_id",
			Message: err.Error(),
		}
	}
	if errors.Is(err, domain.ErrInvalidName) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_name",
			Message: err.Error(),
		}
	}
	if errors.Is(err, domain.ErrInvalidBalance) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_balance",
			Message: err.Error(),
		}
	}

	// Application list errors
	if errors.Is(err, list.ErrCreatedAtWithoutAccountID) || errors.Is(err, list.ErrAccountIDWithoutCreatedAt) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_cursor",
			Message: err.Error(),
		}
	}

	// Application by id errors
	if errors.Is(err, byid.ErrAccountDoesntBelongToUser) {
		return http.StatusForbidden, response.APIError{
			Code:    "account_doesnt_belong_to_user",
			Message: err.Error(),
		}
	}

	// Default
	return http.StatusInternalServerError, response.APIError{
		Code:    "internal_error",
		Message: "internal server error",
	}
}
