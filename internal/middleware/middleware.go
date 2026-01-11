package middleware

import (
	"net/http"
)

type Middleware struct {
}

func New() *Middleware {
	return &Middleware{}
}

func ApplyMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}