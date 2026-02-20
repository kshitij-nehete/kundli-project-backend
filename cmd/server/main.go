package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/kshitij-nehete/astro-report/internal/config"
	"github.com/kshitij-nehete/astro-report/internal/database"
	"github.com/kshitij-nehete/astro-report/internal/server"
	"github.com/kshitij-nehete/astro-report/pkg/logger"
)

func main() {

	cfg := config.LoadConfig()

	mongoClient, err := database.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatal("mongo connection failed", zap.Error(err))
	}

	db := mongoClient.Database(cfg.Database)

	if err := database.CreateUserIndexes(db); err != nil {
		log.Fatal("failed to create indexes", zap.Error(err))
	}

	log, err := logger.NewLogger(cfg.Environment)
	if err != nil {
		panic("failed to initialize logger")
	}
	defer log.Sync()

	httpServer := server.NewHTTPServer(cfg.Port, log, db)

	go func() {
		log.Info("server starting", zap.String("port", cfg.Port))
		if err := httpServer.Start(); err != nil {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", zap.Error(err))
	}

	log.Info("server exited properly")
}
