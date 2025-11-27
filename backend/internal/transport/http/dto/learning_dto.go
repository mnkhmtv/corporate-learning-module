package dto

import "github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"

// UpdateLearningDTO represents full learning update (admin only)
type UpdateLearningDTO struct {
	Topic    string                `json:"topic" binding:"required" example:"Go Programming"`
	Status   string                `json:"status" binding:"required,oneof=active completed" example:"active"`
	Plan     []LearningPlanItemDTO `json:"plan" binding:"required"`
	Feedback *FeedbackDTO          `json:"feedback,omitempty"`
	Notes    *string               `json:"notes,omitempty" example:"Student is making good progress"`
}

// UpdatePlanDTO represents plan update
type UpdatePlanDTO struct {
	Plan []LearningPlanItemDTO `json:"plan" binding:"required"`
}

// LearningPlanItemDTO represents a plan item
type LearningPlanItemDTO struct {
	ID        string `json:"id" binding:"required" example:"1"`
	Text      string `json:"text" binding:"required" example:"Learn Go basics"`
	Completed bool   `json:"completed" example:"false"`
}

// UpdateNotesDTO represents notes update
type UpdateNotesDTO struct {
	Notes string `json:"notes" example:"Student completed module 1"`
}

// CompleteLearningDTO represents learning completion with feedback
type CompleteLearningDTO struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5" example:"5"`
	Comment string `json:"comment" binding:"required" example:"Excellent progress and understanding"`
}

// FeedbackDTO represents feedback
type FeedbackDTO struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5" example:"4"`
	Comment string `json:"comment" binding:"required" example:"Great work!"`
}

// ToPlanItems converts DTOs to domain plan items
func ToPlanItems(dtos []LearningPlanItemDTO) []domain.LearningPlanItem {
	items := make([]domain.LearningPlanItem, len(dtos))
	for i, dto := range dtos {
		items[i] = domain.LearningPlanItem{
			ID:        dto.ID,
			Text:      dto.Text,
			Completed: dto.Completed,
		}
	}
	return items
}

// FromPlanItems converts domain plan items to DTOs
func FromPlanItems(items []domain.LearningPlanItem) []LearningPlanItemDTO {
	dtos := make([]LearningPlanItemDTO, len(items))
	for i, item := range items {
		dtos[i] = LearningPlanItemDTO{
			ID:        item.ID,
			Text:      item.Text,
			Completed: item.Completed,
		}
	}
	return dtos
}
