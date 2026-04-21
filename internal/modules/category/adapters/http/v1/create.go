package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/category/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type createRequest struct {
	Name string `json:"name"`
}

type categoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type createResponse struct {
	Category categoryResponse `json:"category"`
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

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_json",
			Message: "invalid json body",
		}, rid)
		return
	}

	res, err := h.createUC.Execute(r.Context(), create.Command{
		UserID: subID,
		Name:   req.Name,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := createResponse{
		Category: mapCategory(res.Category),
	}

	response.WriteJSON(w, http.StatusCreated, resp, nil, rid)
}

func mapCategory(category *domain.Category) categoryResponse {
	return categoryResponse{
		ID:        category.ID().String(),
		Name:      category.Name().String(),
		UserID:    category.UserID().String(),
		CreatedAt: category.CreatedAt().UTC().Format(time.RFC3339),
		UpdatedAt: category.UpdatedAt().UTC().Format(time.RFC3339),
	}
}
