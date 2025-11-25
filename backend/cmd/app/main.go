package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // драйвер БД
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// для embed
	"github.com/mnkhmtv/corporate-learning-module/backend/config"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/repository/postgres"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	transport "github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http"
	"github.com/mnkhmtv/corporate-learning-module/backend/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		os.Exit(1)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.Env)
	slog.SetDefault(appLogger)

	// Initialize database
	dbConfig := postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	ctx, cancelStart := context.WithCancel(context.Background())
	defer cancelStart()

	pool, err := postgres.NewPool(ctx, dbConfig)
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		log.Fatalf("Database connection failed: %v", err)
		os.Exit(1)
	}
	defer postgres.Close(pool)

	appLogger.Info("Database connection established")

	// Run database migrations
	if cfg.Database.AutoMigrate {
		appLogger.Info("Auto-migration enabled, running migrations...")

		if err := runMigrations(cfg.Database); err != nil {
			appLogger.Error("Migration failed", "error", err)
			log.Fatalf("Migration failed: %v", err)
			os.Exit(1)
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

	ctx, cancelStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStop()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
		log.Fatal("Server shutdown failed:", err)
	}

	appLogger.Info("Server exited successfully")
}

// runMigrations applies all pending migrations from file system
func runMigrations(dbCfg config.DatabaseConfig) error {
	// Construct database URL
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName, dbCfg.SSLMode,
	)

	// Create migrate instance with file:// source
	m, err := migrate.New(
		"file://internal/repository/migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("migration up failed: %w", err)
	}

	log.Println("All migrations applied successfully")
	return nil
}
