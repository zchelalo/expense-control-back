package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/login"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/logout"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/refresh"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/register"
)

type Router struct {
	handler       *Handler
	secureCookies bool
	middleware		*middleware.Middleware
}

func NewRouter(
	registerUC *register.UseCase,
	loginUC *login.UseCase,
	logoutUC *logout.UseCase,
	refreshUC *refresh.UseCase,
	secureCookies bool,
	middleware *middleware.Middleware) *Router {
	handler := NewHandler(
		registerUC,
		loginUC,
		logoutUC,
		refreshUC,
		secureCookies,
	)

	return &Router{
		handler:       handler,
		secureCookies: secureCookies,
		middleware:    middleware,
	}
}

func (r *Router) Register(mux *http.ServeMux) {
	mux.Handle("POST /auth/register", http.HandlerFunc(r.handler.Register))
	mux.Handle("POST /auth/login", http.HandlerFunc(r.handler.Login))
	mux.Handle("POST /auth/logout", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.Logout), r.middleware.Auth))
	mux.Handle("POST /auth/refresh", http.HandlerFunc(r.handler.Refresh))
}
