package v1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type createRequest struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type accountResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	UserID    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type createResponse struct {
	Account accountResponse `json:"account"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_json",
			Message: "invalid json body",
		}, rid)
		return
	}

	if strings.TrimSpace(req.Name) == "" || req.Balance < 0 {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_input",
			Message: "name and balance are required",
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

	balance, err := domain.NewBalance(req.Balance)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_balance",
			Message: "invalid account balance",
		}, rid)
		return
	}

	res, err := h.createUC.Execute(r.Context(), create.Command{
		UserID:  userID,
		Name:    name,
		Balance: balance,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := createResponse{
		Account: accountResponse{
			ID:        res.Account.ID().String(),
			Name:      res.Account.Name().String(),
			Balance:   res.Account.Balance().Float64(),
			UserID:    res.Account.UserID().String(),
			CreatedAt: res.Account.CreatedAt().UTC().Format(time.RFC3339),
			UpdatedAt: res.Account.UpdatedAt().UTC().Format(time.RFC3339),
		},
	}

	response.WriteJSON(w, http.StatusCreated, resp, nil, rid)
}