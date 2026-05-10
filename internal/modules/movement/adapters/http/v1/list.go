package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/list"
	"github.com/zchelalo/expense-control-back/pkg/pagination"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type listResponse struct {
	Movements  []movementResponse `json:"movements"`
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
	language := requestLanguage(r.Header.Get("Accept-Language"))

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

	var createdAt *time.Time
	var movementID *uuid.UUID
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
		movementID = &uid
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
		movementID = &uid
		isBefore = true
	}

	filters, apiErr := parseMovementQueryFilters(queries)
	if apiErr != nil {
		response.WriteError(w, http.StatusBadRequest, *apiErr, rid)
		return
	}

	res, err := h.listUC.Execute(r.Context(), list.Command{
		UserID:         subID,
		AccountID:      filters.AccountID,
		CategoryID:     filters.CategoryID,
		MovementTypeID: filters.MovementTypeID,
		DateFrom:       filters.DateFrom,
		DateTo:         filters.DateTo,
		CreatedAt:      createdAt,
		MovementID:     movementID,
		Limit:          limit,
		IsBefore:       isBefore,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := listResponse{
		Movements: make([]movementResponse, 0, len(res.Movements)),
	}

	for _, movement := range res.Movements {
		resp.Movements = append(resp.Movements, mapMovementDetails(movement, language))
	}

	if len(res.Movements) > 0 {
		lastMovement := res.Movements[len(res.Movements)-1].Movement()
		firstMovement := res.Movements[0].Movement()

		if len(res.Movements) == limit {
			resp.NextCursor = pagination.EncodeCursor(lastMovement.CreatedAt(), lastMovement.ID().UUID())
		}

		if afterCursor != "" || beforeCursor != "" {
			resp.PrevCursor = pagination.EncodeCursor(firstMovement.CreatedAt(), firstMovement.ID().UUID())
		}
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}
