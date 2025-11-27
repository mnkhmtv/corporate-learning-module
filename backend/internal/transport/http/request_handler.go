package http

import (
	"net/http"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/dto"

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

// CreateRequestDTO represents the input for creating a training request
type CreateRequestDTO struct {
	Topic       string `json:"topic" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateRequestDTO represents request update input
type UpdateRequestDTO struct {
	Topic       string `json:"topic" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (h *RequestHandler) CreateRequest(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req CreateRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request, err := h.requestService.CreateRequest(
		c.Request.Context(),
		userID.(string),
		req.Topic,
		req.Description,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTO
	responseDTO := dto.ToRequestResponseDTO(request)
	c.JSON(http.StatusCreated, responseDTO)
}

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

	// Convert to response DTOs
	responseDTOs := dto.ToRequestResponseDTOs(requests)
	c.JSON(http.StatusOK, gin.H{"requests": responseDTOs})
}

func (h *RequestHandler) GetMyRequests(c *gin.Context) {
	userID, _ := c.Get("userID")

	requests, err := h.requestService.GetUserRequests(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTOs
	responseDTOs := dto.ToRequestResponseDTOs(requests)
	c.JSON(http.StatusOK, gin.H{"requests": responseDTOs})
}

func (h *RequestHandler) GetRequestByID(c *gin.Context) {
	requestID := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	request, err := h.requestService.GetRequestByID(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Access control: owner or admin
	if request.UserID != userID.(string) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Convert to response DTO
	responseDTO := dto.ToRequestResponseDTO(request)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *RequestHandler) UpdateRequest(c *gin.Context) {
	requestID := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var req UpdateRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check access: owner or admin
	existingRequest, err := h.requestService.GetRequestByID(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Access control: owner or admin
	if existingRequest.UserID != userID.(string) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	request, err := h.requestService.UpdateRequest(
		c.Request.Context(),
		requestID,
		req.Topic,
		req.Description,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTO
	responseDTO := dto.ToRequestResponseDTO(request)
	c.JSON(http.StatusOK, responseDTO)
}

// AssignMentor handles POST /api/requests/:id/assign (admin only)
func (h *RequestHandler) AssignMentor(c *gin.Context) {
	requestID := c.Param("id")

	var req dto.AssignMentorDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	learning, err := h.requestService.AssignMentor(
		c.Request.Context(),
		requestID,
		req.MentorID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseDTO := dto.ToLearningResponseDTO(learning)
	c.JSON(http.StatusCreated, responseDTO)
}
