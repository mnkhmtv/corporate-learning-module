package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/mnkhmtv/corporate-learning-module/backend/config"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/repository/postgres"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize database connection
	ctx := context.Background()
	dbConfig := postgres.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	pool, err := postgres.NewPool(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	logger.Info("Connected to database successfully")

	// Run database migrations
	if cfg.Database.AutoMigrate {
		logger.Info("Auto-migration enabled, running migrations...")

		if err := runMigrations(cfg.Database); err != nil {
			logger.Error("Migration failed", "error", err)
			log.Fatalf("Migration failed: %v", err)
			os.Exit(1)
		}
		logger.Info("Migrations completed successfully")
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(pool)
	requestRepo := postgres.NewRequestRepository(pool)
	mentorRepo := postgres.NewMentorRepository(pool)
	learningRepo := postgres.NewLearningRepository(pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, cfg.Auth.TokenTTL)
	userService := service.NewUserService(userRepo)
	requestService := service.NewRequestService(requestRepo, userRepo)
	mentorService := service.NewMentorService(mentorRepo)
	learningService := service.NewLearningService(learningRepo, mentorRepo, requestRepo)

	// Initialize HTTP handler
	handler := http.NewHandler(
		authService,
		userService,
		requestService,
		learningService,
		mentorService,
	)

	// Setup Gin router
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	handler.InitRoutes(router, logger, cfg.Auth.JWTSecret)

	// Start server
	addr := ":" + cfg.Server.Port
	logger.Info("Starting server", "address", addr, "environment", cfg.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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
