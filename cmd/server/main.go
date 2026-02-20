package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/kshitij-nehete/astro-report/internal/config"
	"github.com/kshitij-nehete/astro-report/internal/server"
	"github.com/kshitij-nehete/astro-report/pkg/logger"
)

func main() {

	cfg := config.LoadConfig()

	log, err := logger.NewLogger(cfg.Environment)
	if err != nil {
		panic("failed to initialize logger")
	}
	defer log.Sync()

	httpServer := server.NewHTTPServer(cfg.Port, log)

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
