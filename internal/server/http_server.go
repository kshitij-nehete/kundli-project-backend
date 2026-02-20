package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/kshitij-nehete/astro-report/internal/handler"
	"github.com/kshitij-nehete/astro-report/internal/middleware"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(port string, logger *zap.Logger) *HTTPServer {

	r := chi.NewRouter()

	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.LoggingMiddleware(logger))

	r.Get("/health", handler.HealthHandler)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &HTTPServer{
		server: srv,
	}
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
