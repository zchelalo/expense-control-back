package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/pkg/response"
)

type movementTypeResponse struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func mapError(_ error) (int, response.APIError) {
	return http.StatusInternalServerError, response.APIError{
		Code:    "internal_error",
		Message: "internal server error",
	}
}
