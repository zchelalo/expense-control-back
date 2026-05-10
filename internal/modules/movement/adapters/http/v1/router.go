package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/byid"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/delete"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/list"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/stats"
)

type Router struct {
	handler    *Handler
	middleware *middleware.Middleware
}

func NewRouter(
	createUC *create.UseCase,
	listUC *list.UseCase,
	byIDUC *byid.UseCase,
	deleteUC *delete.UseCase,
	statsUC *stats.UseCase,
	middleware *middleware.Middleware,
) *Router {
	handler := NewHandler(createUC, listUC, byIDUC, deleteUC, statsUC)

	return &Router{
		handler:    handler,
		middleware: middleware,
	}
}

func (r *Router) Register(mux *http.ServeMux) {
	mux.Handle("POST /movement/{account_id}", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.Create), r.middleware.Auth))
	mux.Handle("GET /movement", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.List), r.middleware.Auth))
	mux.Handle("GET /movement/stats/overview", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.StatsOverview), r.middleware.Auth))
	mux.Handle("GET /movement/stats/by-account", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.StatsByAccount), r.middleware.Auth))
	mux.Handle("GET /movement/stats/by-category", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.StatsByCategory), r.middleware.Auth))
	mux.Handle("GET /movement/{id}", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.ByID), r.middleware.Auth))
	mux.Handle("DELETE /movement/{id}", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.Delete), r.middleware.Auth))
}
