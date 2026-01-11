package server

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zchelalo/expense-control-back/internal/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Server struct {
	httpServer *http.Server
}

func New(address string, mdw *middleware.Middleware, routerRegistrations ...func(*http.ServeMux)) (*Server, error) {
	mux := http.NewServeMux()

	mux.Handle("GET /metrics", otelhttp.NewHandler(promhttp.Handler(), "GET /metrics"))
	mux.Handle("GET /health", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), "GET /health"))

	for _, register := range routerRegistrations {
		register(mux)
	}

	var handler http.Handler = mux
	handler = otelhttp.NewHandler(handler, "http.server")

	// handler = middleware.ApplyMiddlewares(handler, mdw.XMiddleware...)

	srv := &http.Server{
		Addr:              address,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &Server{httpServer: srv}, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}