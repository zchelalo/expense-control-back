package v1

import (
	"errors"
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/list"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

func mapError(err error) (int, response.APIError) {
	if errors.As(err, &ports.ErrNotFound{}) {
		return http.StatusNotFound, response.APIError{
			Code:    "not_found",
			Message: err.Error(),
		}
	}

	if errors.Is(err, domain.ErrInvalidMovementID) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_movement_id",
			Message: err.Error(),
		}
	}
	if errors.Is(err, domain.ErrInvalidMovementTypeID) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_movement_type_id",
			Message: err.Error(),
		}
	}
	if errors.Is(err, domain.ErrInvalidCategoryID) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_category_id",
			Message: err.Error(),
		}
	}
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
	if errors.Is(err, domain.ErrInvalidDescription) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_description",
			Message: err.Error(),
		}
	}
	if errors.Is(err, domain.ErrInvalidAmount) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_amount",
			Message: err.Error(),
		}
	}
	if errors.Is(err, domain.ErrInsufficientAccountBalance) {
		return http.StatusConflict, response.APIError{
			Code:    "insufficient_account_balance",
			Message: err.Error(),
		}
	}

	if errors.Is(err, list.ErrCreatedAtWithoutMovementID) || errors.Is(err, list.ErrMovementIDWithoutCreatedAt) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_cursor",
			Message: err.Error(),
		}
	}

	return http.StatusInternalServerError, response.APIError{
		Code:    "internal_error",
		Message: "internal server error",
	}
}
