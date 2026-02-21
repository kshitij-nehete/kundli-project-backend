package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/kshitij-nehete/astro-report/internal/handler"
	"github.com/kshitij-nehete/astro-report/internal/middleware"
	"github.com/kshitij-nehete/astro-report/internal/repository"
	"github.com/kshitij-nehete/astro-report/internal/usecase"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(
	port string,
	logger *zap.Logger,
	db *mongo.Database,
) *HTTPServer {

	r := chi.NewRouter()

	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.LoggingMiddleware(logger))

	// Initialize repositories
	userRepo := repository.NewMongoUserRepository(db)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)

	r.Get("/health", handler.HealthHandler(db))
	r.Post("/auth/register", authHandler.Register)

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
