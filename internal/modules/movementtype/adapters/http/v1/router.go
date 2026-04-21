package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movementtype/application/list"
)

type Router struct {
	handler    *Handler
	middleware *middleware.Middleware
}

func NewRouter(
	listUC *list.UseCase,
	middleware *middleware.Middleware,
) *Router {
	return &Router{
		handler:    NewHandler(listUC),
		middleware: middleware,
	}
}

func (r *Router) Register(mux *http.ServeMux) {
	mux.Handle("GET /movement-type", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.List), r.middleware.Auth))
}
