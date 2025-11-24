package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mnkhmtv/corporate-learning-module/backend/config"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/repository/postgres"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	transport "github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http"
	"github.com/mnkhmtv/corporate-learning-module/backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.Env)

	// Initialize database
	dbConfig := postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, dbConfig)
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Database connection failed: %v", err)
	}
	defer postgres.Close(pool)

	appLogger.Info("Database connection established")

	// Auto-migrate if enabled
	if cfg.Database.AutoMigrate {
		appLogger.Info("Auto-migration enabled, running migrations...")
		if err := runMigrations(cfg.Database); err != nil {
			appLogger.Error("Migration failed", "error", err)
			log.Fatalf("Migration failed: %v", err)
		}
		appLogger.Info("Migrations completed successfully")
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(pool)
	mentorRepo := postgres.NewMentorRepository(pool)
	requestRepo := postgres.NewRequestRepository(pool)
	learningRepo := postgres.NewLearningRepository(pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, cfg.Auth.TokenTTL)
	requestService := service.NewRequestService(requestRepo, userRepo)
	mentorService := service.NewMentorService(mentorRepo)
	learningService := service.NewLearningService(learningRepo, requestRepo, mentorRepo, mentorService)

	// Initialize HTTP handler
	handler := transport.NewHandler(authService, requestService, learningService, mentorService)

	// Setup Gin router
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	handler.InitRoutes(router, appLogger, cfg.Auth.JWTSecret)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		appLogger.Info("Starting HTTP server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server failed to start", "error", err)
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
		log.Fatal("Server shutdown failed:", err)
	}

	appLogger.Info("Server exited successfully")
}

func runMigrations(dbCfg config.DatabaseConfig) error {
	// TODO: Implement migration logic using golang-migrate
	// For now, this is a placeholder
	fmt.Println("Migrations would run here")
	return nil
}
