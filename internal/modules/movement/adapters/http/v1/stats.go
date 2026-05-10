package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/stats"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/localization"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

type movementStatsTotalResponse struct {
	Count int64   `json:"count"`
	Total float64 `json:"total"`
}

type movementStatsOverviewResponse struct {
	TotalMovements int64                      `json:"total_movements"`
	Income         movementStatsTotalResponse `json:"income"`
	Expense        movementStatsTotalResponse `json:"expense"`
	NetTotal       float64                    `json:"net_total"`
}

type statsOverviewResponse struct {
	Overview movementStatsOverviewResponse `json:"overview"`
}

type movementStatsAccountResponse struct {
	Account       accountResponse `json:"account"`
	MovementCount int64           `json:"movement_count"`
	IncomeCount   int64           `json:"income_count"`
	ExpenseCount  int64           `json:"expense_count"`
	IncomeTotal   float64         `json:"income_total"`
	ExpenseTotal  float64         `json:"expense_total"`
	NetTotal      float64         `json:"net_total"`
}

type statsByAccountResponse struct {
	Accounts []movementStatsAccountResponse `json:"accounts"`
}

type movementStatsCategoryResponse struct {
	Category      categoryResponse `json:"category"`
	MovementCount int64            `json:"movement_count"`
	IncomeCount   int64            `json:"income_count"`
	ExpenseCount  int64            `json:"expense_count"`
	IncomeTotal   float64          `json:"income_total"`
	ExpenseTotal  float64          `json:"expense_total"`
	NetTotal      float64          `json:"net_total"`
}

type statsByCategoryResponse struct {
	Categories []movementStatsCategoryResponse `json:"categories"`
}

func (h *Handler) StatsOverview(w http.ResponseWriter, r *http.Request) {
	rid := middleware.RequestIDFrom(r.Context())

	subID, ok := middleware.SubjectIDFrom(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, response.APIError{
			Code:    "unauthorized",
			Message: "subject id not found in context",
		}, rid)
		return
	}

	cmd, apiErr := buildStatsCommand(subID, r)
	if apiErr != nil {
		response.WriteError(w, http.StatusBadRequest, *apiErr, rid)
		return
	}

	res, err := h.statsUC.Overview(r.Context(), cmd)
	if err != nil {
		status, mapped := mapError(err)
		response.WriteError(w, status, mapped, rid)
		return
	}

	response.WriteJSON(w, http.StatusOK, statsOverviewResponse{
		Overview: movementStatsOverviewResponse{
			TotalMovements: res.Overview.TotalMovements,
			Income: movementStatsTotalResponse{
				Count: res.Overview.IncomeCount,
				Total: res.Overview.IncomeTotal,
			},
			Expense: movementStatsTotalResponse{
				Count: res.Overview.ExpenseCount,
				Total: res.Overview.ExpenseTotal,
			},
			NetTotal: res.Overview.NetTotal,
		},
	}, nil, rid)
}

func (h *Handler) StatsByAccount(w http.ResponseWriter, r *http.Request) {
	rid := middleware.RequestIDFrom(r.Context())

	subID, ok := middleware.SubjectIDFrom(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, response.APIError{
			Code:    "unauthorized",
			Message: "subject id not found in context",
		}, rid)
		return
	}

	cmd, apiErr := buildStatsCommand(subID, r)
	if apiErr != nil {
		response.WriteError(w, http.StatusBadRequest, *apiErr, rid)
		return
	}

	res, err := h.statsUC.ByAccount(r.Context(), cmd)
	if err != nil {
		status, mapped := mapError(err)
		response.WriteError(w, status, mapped, rid)
		return
	}

	resp := statsByAccountResponse{
		Accounts: make([]movementStatsAccountResponse, 0, len(res.Accounts)),
	}
	for _, account := range res.Accounts {
		resp.Accounts = append(resp.Accounts, mapMovementStatsByAccount(account))
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}

func (h *Handler) StatsByCategory(w http.ResponseWriter, r *http.Request) {
	rid := middleware.RequestIDFrom(r.Context())
	language := requestLanguage(r.Header.Get("Accept-Language"))

	subID, ok := middleware.SubjectIDFrom(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, response.APIError{
			Code:    "unauthorized",
			Message: "subject id not found in context",
		}, rid)
		return
	}

	cmd, apiErr := buildStatsCommand(subID, r)
	if apiErr != nil {
		response.WriteError(w, http.StatusBadRequest, *apiErr, rid)
		return
	}

	res, err := h.statsUC.ByCategory(r.Context(), cmd)
	if err != nil {
		status, mapped := mapError(err)
		response.WriteError(w, status, mapped, rid)
		return
	}

	resp := statsByCategoryResponse{
		Categories: make([]movementStatsCategoryResponse, 0, len(res.Categories)),
	}
	for _, category := range res.Categories {
		resp.Categories = append(resp.Categories, mapMovementStatsByCategory(category, language))
	}

	response.WriteJSON(w, http.StatusOK, resp, nil, rid)
}

func buildStatsCommand(userID uuid.UUID, r *http.Request) (stats.Command, *response.APIError) {
	filters, apiErr := parseMovementQueryFilters(r.URL.Query())
	if apiErr != nil {
		return stats.Command{}, apiErr
	}

	return stats.Command{
		UserID:         userID,
		AccountID:      filters.AccountID,
		CategoryID:     filters.CategoryID,
		MovementTypeID: filters.MovementTypeID,
		DateFrom:       filters.DateFrom,
		DateTo:         filters.DateTo,
	}, nil
}

func mapMovementStatsByAccount(item ports.MovementStatsByAccount) movementStatsAccountResponse {
	return movementStatsAccountResponse{
		Account: accountResponse{
			ID:   item.AccountID.String(),
			Name: item.AccountName,
		},
		MovementCount: item.MovementCount,
		IncomeCount:   item.IncomeCount,
		ExpenseCount:  item.ExpenseCount,
		IncomeTotal:   item.IncomeTotal,
		ExpenseTotal:  item.ExpenseTotal,
		NetTotal:      item.NetTotal,
	}
}

func mapMovementStatsByCategory(item ports.MovementStatsByCategory, language string) movementStatsCategoryResponse {
	name := item.CategoryName
	if item.CategoryIsSystem {
		if localized, ok := localization.LocalizeSystemCategoryName(item.CategorySystemKey, language); ok {
			name = localized
		}
	}

	return movementStatsCategoryResponse{
		Category: categoryResponse{
			ID:        item.CategoryID.String(),
			Name:      name,
			IsSystem:  item.CategoryIsSystem,
			SystemKey: item.CategorySystemKey,
		},
		MovementCount: item.MovementCount,
		IncomeCount:   item.IncomeCount,
		ExpenseCount:  item.ExpenseCount,
		IncomeTotal:   item.IncomeTotal,
		ExpenseTotal:  item.ExpenseTotal,
		NetTotal:      item.NetTotal,
	}
}
