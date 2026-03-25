package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/delete"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type deleteResponse struct {
	Success bool `json:"success"`
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	rid := middleware.RequestIDFrom(r.Context())

	subID, ok := middleware.SubjectIDFrom(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, response.APIError{
			Code:    "unauthorized",
			Message: "subject id not found in context",
		}, rid)
		return
	}

	accountIDString := r.PathValue("id")
	accountIDUUID, err := uuid.Parse(accountIDString)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_account_id",
			Message: "invalid account ID format",
		}, rid)
		return
	}

	res, err := h.deleteUC.Execute(r.Context(), delete.Command{
		UserID:    subID.UUID(),
		AccountID: accountIDUUID,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := deleteResponse{
		Success: res.Success,
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}