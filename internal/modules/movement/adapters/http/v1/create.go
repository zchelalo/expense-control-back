package v1

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	uuidparse "github.com/zchelalo/expense-control-back/pkg/parse"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type createRequest struct {
	Amount         float64 `json:"amount"`
	Description    string  `json:"description"`
	MovementTypeID string  `json:"movement_type_id"`
	CategoryID     string  `json:"category_id"`
}

type movementTypeResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

type categoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type accountResponse struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type movementResponse struct {
	ID           string               `json:"id"`
	Amount       float64              `json:"amount"`
	Description  string               `json:"description"`
	MovementType movementTypeResponse `json:"movement_type"`
	Category     categoryResponse     `json:"category"`
	Account      accountResponse      `json:"account"`
	UserID       string               `json:"user_id"`
	CreatedAt    string               `json:"created_at"`
	UpdatedAt    string               `json:"updated_at"`
}

type createResponse struct {
	Movement movementResponse `json:"movement"`
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

	accountID, err := uuidparse.UUID(r.PathValue("account_id"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_account_id",
			Message: "invalid account ID format",
		}, rid)
		return
	}

	movementTypeID, err := uuidparse.UUID(req.MovementTypeID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_movement_type_id",
			Message: "invalid movement type ID format",
		}, rid)
		return
	}

	categoryID, err := uuidparse.UUID(req.CategoryID)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, response.APIError{
			Code:    "invalid_category_id",
			Message: "invalid category ID format",
		}, rid)
		return
	}

	res, err := h.createUC.Execute(r.Context(), create.Command{
		Amount:         req.Amount,
		Description:    req.Description,
		MovementTypeID: movementTypeID,
		CategoryID:     categoryID,
		AccountID:      accountID,
		UserID:         subID,
	})
	if err != nil {
		status, apiErr := mapError(err)
		response.WriteError(w, status, apiErr, rid)
		return
	}

	resp := createResponse{
		Movement: mapMovement(res.Movement),
	}

	response.WriteJSON(w, http.StatusCreated, resp, nil, rid)
}

func mapMovement(m *domain.Movement) movementResponse {
	return movementResponse{
		ID:          m.ID().String(),
		Amount:      m.Amount().Float64(),
		Description: m.Description().String(),
		MovementType: movementTypeResponse{
			ID: m.MovementTypeID().String(),
		},
		Category: categoryResponse{
			ID: m.CategoryID().String(),
		},
		Account: accountResponse{
			ID: m.AccountID().String(),
		},
		UserID:    m.UserID().String(),
		CreatedAt: m.CreatedAt().UTC().Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt().UTC().Format(time.RFC3339),
	}
}
