package http

import (
	"net/http"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/dto"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	requestService  *service.RequestService
	learningService *service.LearningService
}

func NewRequestHandler(requestService *service.RequestService, learningService *service.LearningService) *RequestHandler {
	return &RequestHandler{
		requestService:  requestService,
		learningService: learningService,
	}
}

// CreateRequest handles POST /api/requests
func (h *RequestHandler) CreateRequest(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.CreateRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request, err := h.requestService.CreateRequest(c.Request.Context(), userID, req.Topic, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// GetAllRequests handles GET /api/requests (admin only)
func (h *RequestHandler) GetAllRequests(c *gin.Context) {
	status := c.Query("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	requests, err := h.requestService.GetAllRequests(c.Request.Context(), statusPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

// GetMyRequests handles GET /api/requests/my
func (h *RequestHandler) GetMyRequests(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	requests, err := h.requestService.GetUserRequests(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

// AssignMentor handles POST /api/requests/:id/assign (admin only)
func (h *RequestHandler) AssignMentor(c *gin.Context) {
	requestID := c.Param("id")

	var req dto.AssignMentorDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Approve the request first
	if err := h.requestService.ApproveRequest(c.Request.Context(), requestID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign mentor and create learning process
	learning, err := h.learningService.AssignMentor(c.Request.Context(), requestID, req.MentorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, learning)
}
