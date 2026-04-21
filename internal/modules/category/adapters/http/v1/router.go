package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/category/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/category/application/delete"
	"github.com/zchelalo/expense-control-back/internal/modules/category/application/list"
)

type Router struct {
	handler    *Handler
	middleware *middleware.Middleware
}

func NewRouter(
	createUC *create.UseCase,
	deleteUC *delete.UseCase,
	listUC *list.UseCase,
	middleware *middleware.Middleware,
) *Router {
	return &Router{
		handler:    NewHandler(createUC, deleteUC, listUC),
		middleware: middleware,
	}
}

func (r *Router) Register(mux *http.ServeMux) {
	mux.Handle("POST /category", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.Create), r.middleware.Auth))
	mux.Handle("GET /category", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.List), r.middleware.Auth))
	mux.Handle("DELETE /category/{id}", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.Delete), r.middleware.Auth))
}
