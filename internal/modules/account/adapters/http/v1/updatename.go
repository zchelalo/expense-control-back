package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/updatename"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type updateNameRequest struct {
	Name    string  `json:"name"`
}

type updateNameResponse struct {
	Account   accountResponse `json:"account"`
}

func (h *Handler) UpdateName(w http.ResponseWriter, r *http.Request) {
	rid := middleware.RequestIDFrom(r.Context())

	subID, ok := middleware.SubjectIDFrom(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, response.APIError{
			Code:    "unauthorized",
			Message: "subject id not found in context",
		}, rid)
		return
	}

	userID, err := domain.NewUserID(subID.UUID())
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_user_id",
			Message: "invalid user ID format",
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
	accountID, err := domain.NewAccountID(accountIDUUID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_account_id",
			Message: "invalid account ID format",
		}, rid)
		return
	}

	var req updateNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_json",
			Message: "invalid json body",
		}, rid)
		return
	}

	name, err := domain.NewName(req.Name)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_name",
			Message: "invalid account name",
		}, rid)
		return
	}

	res, err := h.updateNameUC.Execute(r.Context(), updatename.Command{
		UserID:    userID,
		AccountID: accountID,
		Name: name,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := updateNameResponse{
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