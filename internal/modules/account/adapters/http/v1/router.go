package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/byid"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/list"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/updatename"
)

type Router struct {
	handler    *Handler
	middleware *middleware.Middleware
}

func NewRouter(
	createUC *create.UseCase,
	listUC *list.UseCase,
	byIDUC *byid.UseCase,
	updateName *updatename.UseCase,
	middleware *middleware.Middleware,
) *Router {
	handler := NewHandler(
		createUC,
		listUC,
		byIDUC,
		updateName,
	)

	return &Router{
		handler:    handler,
		middleware: middleware,
	}
}

func (r *Router) Register(mux *http.ServeMux) {
	mux.Handle("POST /account", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.Create), r.middleware.Auth))
	mux.Handle("GET /account", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.List), r.middleware.Auth))
	mux.Handle("GET /account/{id}", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.ByID), r.middleware.Auth))
	mux.Handle("PATCH /account/{id}/name", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.UpdateName), r.middleware.Auth))
}
