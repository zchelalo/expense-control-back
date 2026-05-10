package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/byid"
	uuidparse "github.com/zchelalo/expense-control-back/pkg/parse"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type byIDResponse struct {
	Movement movementResponse `json:"movement"`
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

	movementID, err := uuidparse.UUID(r.PathValue("id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_movement_id",
			Message: "invalid movement ID format",
		}, rid)
		return
	}

	language := requestLanguage(r.Header.Get("Accept-Language"))

	res, err := h.byIDUC.Execute(r.Context(), byid.Command{
		UserID:     subID,
		MovementID: movementID,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	response.WriteJSON(w, http.StatusOK, byIDResponse{
		Movement: mapMovementDetails(res.Movement, language),
	}, nil, rid)
}
