package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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
	apiMux := http.NewServeMux()

	apiMux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	apiMux.Handle("GET /metrics", promhttp.Handler())

	for _, register := range routerRegistrations {
		register(apiMux)
	}

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	var handler http.Handler = mux
	handler = otelhttp.NewHandler(handler, "http.server",
		otelhttp.WithFilter(func(r *http.Request) bool {
			return !middleware.ShouldSkipObservability(r.URL.Path)
		}),
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			path := r.URL.Path
			if path != "/" {
				path = strings.TrimSuffix(path, "/")
			}
			return fmt.Sprintf("%s %s", r.Method, path)
		}),
	)

	handler = middleware.ApplyMiddlewares(handler,
		mdw.RequestID,
		mdw.InjectLogger,
		mdw.LogRequest,
		mdw.AccessControl,
	)

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

func (s *Server) Start() error { return s.httpServer.ListenAndServe() }
func (s *Server) Shutdown(ctx context.Context) error { return s.httpServer.Shutdown(ctx) }