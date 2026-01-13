package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/register"
)

type Router struct {
	handler       *Handler
	secureCookies bool
}

func NewRouter(registerUC *register.UseCase, secureCookies bool) *Router {
	handler := NewHandler(registerUC, secureCookies)

	return &Router{
		handler:       handler,
		secureCookies: secureCookies,
	}
}

func (r *Router) Register(mux *http.ServeMux) {
	mux.Handle("POST /v1/auth/register", http.HandlerFunc(r.handler.Register))

	// mux.Handle("POST /v1/auth/login", http.HandlerFunc(r.handler.Login))
	// mux.Handle("POST /v1/auth/refresh", http.HandlerFunc(r.handler.Refresh))
	// mux.Handle("POST /v1/auth/logout", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.Logout), r.middleware.Auth))
}
