package v1

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/byid"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type byIDResponse struct {
	Account   accountResponse `json:"account"`
}

func (h *Handler) ByID(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.byIDUC.Execute(r.Context(), byid.Command{
		UserID:    subID,
		AccountID: accountIDUUID,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := byIDResponse{
		Account: accountResponse{
			ID:        res.Account.ID().String(),
			Name:      res.Account.Name().String(),
			Balance:   res.Account.Balance().Float64(),
			UserID:    res.Account.UserID().String(),
			CreatedAt: res.Account.CreatedAt().UTC().Format(time.RFC3339),
			UpdatedAt: res.Account.UpdatedAt().UTC().Format(time.RFC3339),
		},
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}