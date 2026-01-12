package v1

import (
	"net/http"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/register"
)

type Router struct {
	mux     *http.ServeMux
	handler *Handler
}

func NewRouter(mux *http.ServeMux, registerUC *register.UseCase, secureCookies bool) *Router {
	return &Router{
		mux:     mux,
		handler: NewHandler(registerUC, secureCookies),
	}
}

func (r *Router) SetRoutes() {
	r.mux.Handle("POST /auth/register", http.HandlerFunc(r.handler.Register))

	// r.mux.Handle("POST /auth/login", http.HandlerFunc(r.handler.Login))
	// r.mux.Handle("POST /auth/refresh", http.HandlerFunc(r.handler.Refresh))
	// r.mux.Handle("POST /auth/logout", middleware.ApplyMiddlewares(http.HandlerFunc(r.handler.Logout), r.middleware.Auth))
}