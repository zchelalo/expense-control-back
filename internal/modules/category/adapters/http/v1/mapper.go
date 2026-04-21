package v1

import (
	"errors"
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/modules/category/application/list"
	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/category/ports"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

func mapError(err error) (int, response.APIError) {
	if errors.As(err, &ports.ErrAlreadyExists{}) {
		return http.StatusConflict, response.APIError{
			Code:    "already_exists",
			Message: err.Error(),
		}
	}

	if errors.As(err, &ports.ErrNotFound{}) {
		return http.StatusNotFound, response.APIError{
			Code:    "not_found",
			Message: err.Error(),
		}
	}
	if errors.As(err, &ports.ErrInUse{}) {
		return http.StatusConflict, response.APIError{
			Code:    "in_use",
			Message: err.Error(),
		}
	}

	if errors.Is(err, domain.ErrInvalidCategoryID) {
		return http.StatusBadRequest, response.APIError{
			Code:    "invalid_category_id",
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

	if errors.Is(err, list.ErrCreatedAtWithoutCategoryID) || errors.Is(err, list.ErrCategoryIDWithoutCreatedAt) {
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
