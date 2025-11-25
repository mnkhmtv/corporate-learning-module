package http

import (
	"log/slog"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authHandler     *AuthHandler
	requestHandler  *RequestHandler
	learningHandler *LearningHandler
	mentorHandler   *MentorHandler
}

func NewHandler(
	authService *service.AuthService,
	requestService *service.RequestService,
	learningService *service.LearningService,
	mentorService *service.MentorService,
) *Handler {
	return &Handler{
		authHandler:     NewAuthHandler(authService),
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

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check (no auth required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "internal-training-system",
		})
	})

	// API v1 group
	api := router.Group("/api")
	{
		// Authentication routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.authHandler.Register)
			auth.POST("/login", h.authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(jwtSecret), h.authHandler.GetMe)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Training Requests
			requests := protected.Group("/requests")
			{
				requests.POST("", h.requestHandler.CreateRequest)
				requests.GET("/my", h.requestHandler.GetMyRequests)

				// Admin only routes
				admin := requests.Group("")
				admin.Use(middleware.RequireAdmin())
				{
					admin.GET("", h.requestHandler.GetAllRequests)
					admin.POST("/:id/assign", h.requestHandler.AssignMentor)
				}
			}

			// Learning Processes
			learnings := protected.Group("/learnings")
			{
				learnings.GET("", h.learningHandler.GetMyLearnings)
				learnings.GET("/:id", h.learningHandler.GetLearningByID)
				learnings.GET("/:id/progress", h.learningHandler.GetProgress)

				learnings.POST("/:id/plan", h.learningHandler.AddPlanItem)
				learnings.PUT("/:id/plan", h.learningHandler.UpdatePlan)
				learnings.PUT("/:id/plan/:itemId", h.learningHandler.UpdatePlanItem)
				learnings.PATCH("/:id/plan/:itemId/toggle", h.learningHandler.TogglePlanItem)
				learnings.DELETE("/:id/plan/:itemId", h.learningHandler.RemovePlanItem)

				learnings.PUT("/:id/notes", h.learningHandler.UpdateNotes)

				learnings.POST("/:id/complete", h.learningHandler.CompleteLearning)
			}

			// Mentors
			mentors := protected.Group("/mentors")
			{
				mentors.GET("", h.mentorHandler.GetAllMentors)
				mentors.GET("/:id", h.mentorHandler.GetMentorByID)

				// Admin only
				adminMentors := mentors.Group("")
				adminMentors.Use(middleware.RequireAdmin())
				{
					adminMentors.POST("", h.mentorHandler.CreateMentor)
				}
			}
		}
	}
}
