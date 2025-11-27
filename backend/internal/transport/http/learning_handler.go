package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/dto"
)

type LearningHandler struct {
	learningService *service.LearningService
}

func NewLearningHandler(learningService *service.LearningService) *LearningHandler {
	return &LearningHandler{
		learningService: learningService,
	}
}

func (h *LearningHandler) GetAllLearnings(c *gin.Context) {
	learnings, err := h.learningService.GetAllLearnings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTOs
	responseDTOs := dto.ToLearningResponseDTOs(learnings)
	c.JSON(http.StatusOK, gin.H{"learnings": responseDTOs})
}

func (h *LearningHandler) GetMyLearnings(c *gin.Context) {
	userID, _ := c.Get("userID")

	learnings, err := h.learningService.GetUserLearnings(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTOs
	responseDTOs := dto.ToLearningResponseDTOs(learnings)
	c.JSON(http.StatusOK, gin.H{"learnings": responseDTOs})
}

func (h *LearningHandler) GetLearningByID(c *gin.Context) {
	learningID := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	learning, err := h.learningService.GetLearningByID(c.Request.Context(), learningID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Access control: owner or admin
	if learning.UserID != userID.(string) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Convert to response DTO
	responseDTO := dto.ToLearningResponseDTO(learning)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *LearningHandler) UpdateLearning(c *gin.Context) {
	learningID := c.Param("id")

	var req dto.UpdateLearningDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert DTO to domain
	status := domain.LearningStatus(req.Status)
	plan := dto.ToPlanItems(req.Plan)

	var feedback *domain.Feedback
	if req.Feedback != nil {
		feedback = &domain.Feedback{
			Rating:  req.Feedback.Rating,
			Comment: req.Feedback.Comment,
		}
	}

	learning, err := h.learningService.UpdateLearning(
		c.Request.Context(),
		learningID,
		req.Topic,
		status,
		plan,
		feedback,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTO
	responseDTO := dto.ToLearningResponseDTO(learning)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *LearningHandler) UpdatePlan(c *gin.Context) {
	learningID := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var req dto.UpdatePlanDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check access: owner or admin
	existingLearning, err := h.learningService.GetLearningByID(c.Request.Context(), learningID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if existingLearning.UserID != userID.(string) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	plan := dto.ToPlanItems(req.Plan)

	learning, err := h.learningService.UpdatePlan(c.Request.Context(), learningID, plan)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTO
	responseDTO := dto.ToLearningResponseDTO(learning)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *LearningHandler) UpdateNotes(c *gin.Context) {
	learningID := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var req dto.UpdateNotesDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check access: owner or admin
	existingLearning, err := h.learningService.GetLearningByID(c.Request.Context(), learningID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if existingLearning.UserID != userID.(string) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	learning, err := h.learningService.UpdateNotes(c.Request.Context(), learningID, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTO
	responseDTO := dto.ToLearningResponseDTO(learning)
	c.JSON(http.StatusOK, responseDTO)
}

func (h *LearningHandler) CompleteLearning(c *gin.Context) {
	learningID := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var req dto.CompleteLearningDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check access: owner or admin
	existingLearning, err := h.learningService.GetLearningByID(c.Request.Context(), learningID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if existingLearning.UserID != userID.(string) && role.(string) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	learning, err := h.learningService.CompleteLearning(
		c.Request.Context(),
		learningID,
		req.Rating,
		req.Comment,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to response DTO
	responseDTO := dto.ToLearningResponseDTO(learning)
	c.JSON(http.StatusOK, responseDTO)
}
