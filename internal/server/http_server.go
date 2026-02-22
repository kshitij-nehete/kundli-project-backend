package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/kshitij-nehete/astro-report/internal/auth"
	"github.com/kshitij-nehete/astro-report/internal/config"
	"github.com/kshitij-nehete/astro-report/internal/handler"
	"github.com/kshitij-nehete/astro-report/internal/middleware"
	"github.com/kshitij-nehete/astro-report/internal/repository"
	"github.com/kshitij-nehete/astro-report/internal/response"
	"github.com/kshitij-nehete/astro-report/internal/usecase"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(
	port string,
	logger *zap.Logger,
	db *mongo.Database,
	cfg *config.Config,
) *HTTPServer {

	r := chi.NewRouter()

	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.SecurityHeadersMiddleware)
	r.Use(middleware.RecoveryMiddleware(logger))
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.RateLimitMiddleware)

	// Initialize repositories
	userRepo := repository.NewMongoUserRepository(db)

	reportRepo := repository.NewMongoReportRepository(db)
	reportUsecase := usecase.NewReportUsecase(reportRepo, userRepo)
	reportHandler := handler.NewReportHandler(reportUsecase)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo)

	// Initialize JWT service
	jwtService := auth.NewJWTService(cfg.JWTSecret)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase, jwtService)

	r.Get("/health", handler.HealthHandler(db))
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware(jwtService))

		r.Get("/auth/me", func(w http.ResponseWriter, r *http.Request) {

			userID, ok := r.Context().Value(middleware.UserIDKey).(string)
			if !ok {
				response.WriteJSONError(w, http.StatusUnauthorized, "invalid context user")
				return
			}

			w.Write([]byte("Authenticated user ID: " + userID))
		})

		r.Post("/reports", reportHandler.Create)
	})

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
