package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movementtype/application/list"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type listResponse struct {
	MovementTypes []movementTypeResponse `json:"movement_types"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	rid := middleware.RequestIDFrom(r.Context())

	res, err := h.listUC.Execute(r.Context(), list.Command{})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := listResponse{
		MovementTypes: make([]movementTypeResponse, 0, len(res.MovementTypes)),
	}

	for _, movementType := range res.MovementTypes {
		resp.MovementTypes = append(resp.MovementTypes, movementTypeResponse{
			ID:          movementType.ID().String(),
			Key:         movementType.Key(),
			Name:        movementType.Name(),
			Description: movementType.Description(),
		})
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}
