package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/create"
)

type Router struct {
	handler       *Handler
	middleware		*middleware.Middleware
}

func NewRouter(
	createUC *create.UseCase,
	middleware *middleware.Middleware) *Router {
	handler := NewHandler(
		createUC,
	)

	return &Router{
		handler:       handler,
		middleware:    middleware,
	}
}

func (r *Router) Register(mux *http.ServeMux) {
	mux.Handle("POST /account/create", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.Create), r.middleware.Auth))
}
