package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/list"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	"github.com/zchelalo/expense-control-back/pkg/pagination"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type listResponse struct {
	Accounts   []accountResponse `json:"accounts"`
	NextCursor string            `json:"next_cursor,omitempty"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
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

	queries := r.URL.Query()

	var limit int
	limitRaw := queries.Get("limit")
	if limitRaw != "" {
		limit, err = strconv.Atoi(limitRaw)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, response.APIError{
				Code:    "invalid_limit",
				Message: "limit must be a valid integer",
			}, rid)
			return
		}
	}

	var createdAt *time.Time
	var accountID *domain.AccountID

	cursorRaw := queries.Get("cursor")
	if cursorRaw != "" {
		ts, uid, err := pagination.DecodeCursor(cursorRaw)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, response.APIError{
				Code:    "invalid_cursor",
				Message: "invalid pagination cursor",
			}, rid)
			return
		}
		createdAt = &ts
		accID, err := domain.NewAccountID(uid)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, response.APIError{
				Code:    "invalid_cursor",
				Message: "invalid account ID in cursor",
			}, rid)
			return
		}
		accountID = &accID
	}

	res, err := h.listUC.Execute(r.Context(), list.Command{
		UserID:    userID,
		CreatedAt: createdAt,
		AccountID: accountID,
		Limit:     limit,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := listResponse{
		Accounts: make([]accountResponse, 0, len(res.Account)),
	}

	for _, acc := range res.Account {
		resp.Accounts = append(resp.Accounts, accountResponse{
			ID:        acc.ID().String(),
			Name:      acc.Name().String(),
			Balance:   acc.Balance().Float64(),
			UserID:    acc.UserID().String(),
			CreatedAt: acc.CreatedAt().UTC().Format(time.RFC3339),
			UpdatedAt: acc.UpdatedAt().UTC().Format(time.RFC3339),
		})
	}

	if len(res.Account) > 0 && len(res.Account) == limit {
		lastAcc := res.Account[len(res.Account)-1]
		resp.NextCursor = pagination.EncodeCursor(lastAcc.CreatedAt(), lastAcc.ID().UUID())
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}