package dto

import "github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"

// UpdateLearningDTO represents full learning update (admin only)
type UpdateLearningDTO struct {
	Topic    string                `json:"topic" binding:"required"`
	Status   string                `json:"status" binding:"required,oneof=active completed"`
	Plan     []LearningPlanItemDTO `json:"plan" binding:"required"`
	Feedback *FeedbackDTO          `json:"feedback,omitempty"`
	Notes    *string               `json:"notes,omitempty"`
}

// UpdatePlanDTO represents plan update
type UpdatePlanDTO struct {
	Plan []LearningPlanItemDTO `json:"plan" binding:"required"`
}

// LearningPlanItemDTO represents a plan item
type LearningPlanItemDTO struct {
	ID        string `json:"id" binding:"required"`
	Text      string `json:"text" binding:"required"`
	Completed bool   `json:"completed"`
}

// UpdateNotesDTO represents notes update
type UpdateNotesDTO struct {
	Notes string `json:"notes"`
}

// CompleteLearningDTO represents learning completion with feedback
type CompleteLearningDTO struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required"`
}

// FeedbackDTO represents feedback
type FeedbackDTO struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required"`
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
