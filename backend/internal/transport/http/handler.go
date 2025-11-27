package http

import (
	"log/slog"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	authHandler     *AuthHandler
	userHandler     *UserHandler
	requestHandler  *RequestHandler
	learningHandler *LearningHandler
	mentorHandler   *MentorHandler
}

func NewHandler(
	authService *service.AuthService,
	userService *service.UserService,
	requestService *service.RequestService,
	learningService *service.LearningService,
	mentorService *service.MentorService,
) *Handler {
	return &Handler{
		authHandler:     NewAuthHandler(authService, userService),
		userHandler:     NewUserHandler(userService, learningService, requestService),
		requestHandler:  NewRequestHandler(requestService, learningService),
		learningHandler: NewLearningHandler(learningService),
		mentorHandler:   NewMentorHandler(mentorService),
	}
}

// InitRoutes registers all HTTP routes
func (h *Handler) InitRoutes(router *gin.Engine, logger *slog.Logger, jwtSecret string) {
	// Global middleware
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.PrometheusMiddleware())

	// Root level endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "internal-training-system",
		})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api")
	{
		// Auth /api/auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.authHandler.Register)
			auth.POST("/login", h.authHandler.Login)

			// Protected auth routes
			auth.GET("/me", middleware.AuthMiddleware(jwtSecret), h.authHandler.GetMe)
			auth.PUT("/me", middleware.AuthMiddleware(jwtSecret), h.authHandler.UpdateMe)
		}

		// Protected routes (require authentication)
		// Users /api/users
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Admin only
			users.GET("", middleware.AdminOnly(), h.userHandler.GetAllUsers)
			users.GET("/:id", middleware.OwnerOrAdminOnly(), h.userHandler.GetUserByID)
			users.PUT("/:id", middleware.OwnerOrAdminOnly(), h.userHandler.UpdateUserByID)
			users.GET("/:id/requests", middleware.OwnerOrAdminOnly(), h.userHandler.GetUserRequests)
			users.GET("/:id/learnings", middleware.OwnerOrAdminOnly(), h.userHandler.GetUserLearnings)
		}

		// Requests /api/requests
		requests := api.Group("/requests")
		requests.Use(middleware.AuthMiddleware(jwtSecret))
		{
			requests.GET("", middleware.AdminOnly(), h.requestHandler.GetAllRequests)
			requests.POST("", h.requestHandler.CreateRequest)
			requests.GET("/my", h.requestHandler.GetMyRequests)
			requests.GET("/:id", h.requestHandler.GetRequestByID)
			requests.PUT("/:id", h.requestHandler.UpdateRequest)
			requests.POST("/:id/assign", middleware.AdminOnly(), h.requestHandler.AssignMentor)
		}

		// Mentors /api/mentors
		mentors := api.Group("/mentors")
		mentors.Use(middleware.AuthMiddleware(jwtSecret))
		{
			mentors.GET("", h.mentorHandler.GetAllMentors)
			mentors.POST("", middleware.AdminOnly(), h.mentorHandler.CreateMentor)
			mentors.GET("/:id", h.mentorHandler.GetMentorByID)
			mentors.PUT("/:id", middleware.AdminOnly(), h.mentorHandler.UpdateMentor)
		}

		// Learnings /api/learnings
		learnings := api.Group("/learnings")
		learnings.Use(middleware.AuthMiddleware(jwtSecret))
		{
			learnings.GET("", h.learningHandler.GetMyLearnings)
			learnings.POST("", h.learningHandler.CreateLearning)
			learnings.GET("/:id", h.learningHandler.GetLearningByID)
			learnings.PUT("/:id", middleware.AdminOnly(), h.learningHandler.UpdateLearning)
			learnings.PUT("/:id/plan", h.learningHandler.UpdatePlan)
			learnings.PUT("/:id/notes", h.learningHandler.UpdateNotes)
			learnings.POST("/:id/complete", h.learningHandler.CompleteLearning)
		}
	}
}
