package dto

import (
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
)

// ToRequestResponseDTO converts domain TrainingRequest to response DTO
func ToRequestResponseDTO(req *domain.TrainingRequest) TrainingRequestResponseDTO {
	return TrainingRequestResponseDTO{
		ID: req.ID,
		User: RequestUserDTO{
			ID:       req.UserID,
			Name:     req.UserName,
			JobTitle: req.UserJobTitle,
			Telegram: req.UserTelegram,
		},
		Topic:       req.Topic,
		Description: req.Description,
		Status:      string(req.Status),
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}
}

// ToRequestResponseDTOs converts slice of domain TrainingRequests to response DTOs
func ToRequestResponseDTOs(requests []*domain.TrainingRequest) []TrainingRequestResponseDTO {
	dtos := make([]TrainingRequestResponseDTO, len(requests))
	for i, req := range requests {
		dtos[i] = ToRequestResponseDTO(req)
	}
	return dtos
}

// ToLearningResponseDTO converts domain LearningProcess to response DTO
func ToLearningResponseDTO(learning *domain.LearningProcess) LearningProcessResponseDTO {
	// Convert plan items
	planItems := make([]LearningPlanItemDTO, len(learning.Plan))
	for i, item := range learning.Plan {
		planItems[i] = LearningPlanItemDTO{
			ID:        item.ID,
			Text:      item.Text,
			Completed: item.Completed,
		}
	}

	// Convert feedback
	var feedbackDTO *FeedbackDTO
	if learning.Feedback != nil {
		feedbackDTO = &FeedbackDTO{
			Rating:  learning.Feedback.Rating,
			Comment: learning.Feedback.Comment,
		}
	}

	return LearningProcessResponseDTO{
		ID: learning.ID,
		Request: LearningRequestDTO{
			ID:          learning.RequestID,
			Topic:       learning.RequestTopic,
			Description: learning.RequestDescription,
		},
		User: LearningUserDTO{
			ID:   learning.UserID,
			Name: learning.UserName,
		},
		Mentor: LearningMentorDTO{
			ID:         learning.MentorID,
			Name:       learning.MentorName,
			Telegram:   learning.MentorTelegram,
			JobTitle:   learning.MentorJobTitle,
			Experience: learning.MentorExperience,
		},
		Status:    string(learning.Status),
		StartDate: learning.StartDate,
		EndDate:   learning.EndDate,
		Plan:      planItems,
		Feedback:  feedbackDTO,
		Notes:     learning.Notes,
	}
}

// ToLearningResponseDTOs converts slice of domain LearningProcesses to response DTOs
func ToLearningResponseDTOs(learnings []*domain.LearningProcess) []LearningProcessResponseDTO {
	dtos := make([]LearningProcessResponseDTO, len(learnings))
	for i, learning := range learnings {
		dtos[i] = ToLearningResponseDTO(learning)
	}
	return dtos
}
