package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/category/application/list"
	"github.com/zchelalo/expense-control-back/pkg/pagination"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type listResponse struct {
	Categories []categoryResponse `json:"categories"`
	PrevCursor string             `json:"prev_cursor,omitempty"`
	NextCursor string             `json:"next_cursor,omitempty"`
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
	if limitRaw := queries.Get("limit"); limitRaw != "" {
		limit, err = strconv.Atoi(limitRaw)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, response.APIError{
				Code:    "invalid_limit",
				Message: "limit must be a valid integer",
			}, rid)
			return
		}
	}

	var (
		createdAt  *time.Time
		categoryID *uuid.UUID
		isBefore   bool
	)

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
		categoryID = &uid
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
		categoryID = &uid
		isBefore = true
	}

	res, err := h.listUC.Execute(r.Context(), list.Command{
		UserID:     subID,
		CreatedAt:  createdAt,
		CategoryID: categoryID,
		Limit:      limit,
		IsBefore:   isBefore,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := listResponse{
		Categories: make([]categoryResponse, 0, len(res.Categories)),
	}

	for i := range res.Categories {
		resp.Categories = append(resp.Categories, mapCategory(&res.Categories[i]))
	}

	if len(res.Categories) > 0 {
		lastCategory := res.Categories[len(res.Categories)-1]
		firstCategory := res.Categories[0]

		if len(res.Categories) == limit {
			resp.NextCursor = pagination.EncodeCursor(lastCategory.CreatedAt(), lastCategory.ID().UUID())
		}

		if afterCursor != "" || beforeCursor != "" {
			resp.PrevCursor = pagination.EncodeCursor(firstCategory.CreatedAt(), firstCategory.ID().UUID())
		}
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}
