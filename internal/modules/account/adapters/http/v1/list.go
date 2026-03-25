package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/list"
	"github.com/zchelalo/expense-control-back/pkg/pagination"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type listResponse struct {
	Accounts   []accountResponse `json:"accounts"`
	PrevCursor string            `json:"prev_cursor,omitempty"`
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

	queries := r.URL.Query()

	var err error
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
	var accountID *uuid.UUID
	var isBefore bool

	afterCursor := queries.Get("after_cursor")
	beforeCursor := queries.Get("before_cursor")

	if afterCursor != "" {
		ts, uid, err := pagination.DecodeCursor(afterCursor)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, response.APIError{
				Code:    "invalid_after_cursor",
				Message: "invalid after pagination cursor",
			}, rid)
			return
		}
		createdAt = &ts
		accountID = &uid
	} else if beforeCursor != "" {
		ts, uid, err := pagination.DecodeCursor(beforeCursor)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, response.APIError{
				Code:    "invalid_before_cursor",
				Message: "invalid before pagination cursor",
			}, rid)
			return
		}
		createdAt = &ts
		accountID = &uid
		isBefore = true
	}

	search := queries.Get("search")
	var name *string
	if search != "" {
		name = &search
	}

	res, err := h.listUC.Execute(r.Context(), list.Command{
		UserID:    subID,
		Name:      name,
		CreatedAt: createdAt,
		AccountID: accountID,
		Limit:     limit,
		IsBefore:  isBefore,
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

	if len(res.Account) > 0 {
		lastAcc := res.Account[len(res.Account)-1]
		firstAcc := res.Account[0]

		if len(res.Account) == limit {
			resp.NextCursor = pagination.EncodeCursor(lastAcc.CreatedAt(), lastAcc.ID().UUID())
		}

		if afterCursor != "" || beforeCursor != "" {
			resp.PrevCursor = pagination.EncodeCursor(firstAcc.CreatedAt(), firstAcc.ID().UUID())
		}
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}