package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/service"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/dto"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/transport/http/middleware"
)

type LearningHandler struct {
	learningService *service.LearningService
}

func NewLearningHandler(learningService *service.LearningService) *LearningHandler {
	return &LearningHandler{
		learningService: learningService,
	}
}

// GetMyLearnings handles GET /api/learnings
func (h *LearningHandler) GetMyLearnings(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	learnings, err := h.learningService.GetUserLearnings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, learnings)
}

// GetLearningByID handles GET /api/learnings/:id
func (h *LearningHandler) GetLearningByID(c *gin.Context) {
	learningID := c.Param("id")

	learning, err := h.learningService.GetLearningByID(c.Request.Context(), learningID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, learning)
}

// AddPlanItem handles POST /api/learnings/:id/plan
func (h *LearningHandler) AddPlanItem(c *gin.Context) {
	learningID := c.Param("id")

	var req dto.AddPlanItemDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	learning, err := h.learningService.AddPlanItem(c.Request.Context(), learningID, req.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, learning)
}

// UpdatePlan handles PUT /api/learnings/:id/plan (обновление всего плана)
func (h *LearningHandler) UpdatePlan(c *gin.Context) {
	learningID := c.Param("id")

	var req dto.UpdatePlanDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Конвертируем DTO в domain модель
	planItems := make([]domain.LearningPlanItem, len(req.Plan))
	for i, item := range req.Plan {
		planItems[i] = domain.LearningPlanItem{
			ID:        item.ID,
			Text:      item.Text,
			Completed: item.Completed,
		}
	}

	learning, err := h.learningService.UpdatePlan(c.Request.Context(), learningID, planItems)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, learning)
}

// UpdatePlanItem handles PUT /api/learnings/:id/plan/:itemId
func (h *LearningHandler) UpdatePlanItem(c *gin.Context) {
	learningID := c.Param("id")
	itemID := c.Param("itemId")

	var req dto.UpdatePlanItemDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.learningService.UpdatePlanItem(c.Request.Context(), learningID, itemID, req.Text, req.Completed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "plan item updated successfully"})
}

// TogglePlanItem handles PATCH /api/learnings/:id/plan/:itemId/toggle
func (h *LearningHandler) TogglePlanItem(c *gin.Context) {
	learningID := c.Param("id")
	itemID := c.Param("itemId")

	err := h.learningService.TogglePlanItem(c.Request.Context(), learningID, itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "plan item toggled successfully"})
}

// RemovePlanItem handles DELETE /api/learnings/:id/plan/:itemId
func (h *LearningHandler) RemovePlanItem(c *gin.Context) {
	learningID := c.Param("id")
	itemID := c.Param("itemId")

	err := h.learningService.RemovePlanItem(c.Request.Context(), learningID, itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "plan item removed successfully"})
}

// CompleteLearning handles POST /api/learnings/:id/complete
func (h *LearningHandler) CompleteLearning(c *gin.Context) {
	learningID := c.Param("id")

	var req dto.CompleteLearningDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.learningService.CompleteLearning(c.Request.Context(), learningID, req.Rating, req.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "learning completed successfully"})
}

// GetProgress handles GET /api/learnings/:id/progress
func (h *LearningHandler) GetProgress(c *gin.Context) {
	learningID := c.Param("id")

	progress, err := h.learningService.GetLearningProgress(c.Request.Context(), learningID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"progress": progress})
}

// UpdateNotes handles PATCH /api/learnings/:id/notes
func (h *LearningHandler) UpdateNotes(c *gin.Context) {
	learningID := c.Param("id")

	var req dto.UpdateNotesDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	learning, err := h.learningService.UpdateNotes(c.Request.Context(), learningID, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, learning)
}
