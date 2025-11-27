package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
)

type LearningService struct {
	learningRepo domain.LearningRepository
	mentorRepo   domain.MentorRepository
	requestRepo  domain.RequestRepository
}

func NewLearningService(
	learningRepo domain.LearningRepository,
	mentorRepo domain.MentorRepository,
	requestRepo domain.RequestRepository,
) *LearningService {
	return &LearningService{
		learningRepo: learningRepo,
		mentorRepo:   mentorRepo,
		requestRepo:  requestRepo,
	}
}

// GetAllLearnings retrieves all learning processes (admin only)
func (s *LearningService) GetAllLearnings(ctx context.Context) ([]*domain.LearningProcess, error) {
	return s.learningRepo.GetAll(ctx)
}

// GetUserLearnings retrieves all learning processes for current user
func (s *LearningService) GetUserLearnings(ctx context.Context, userID string) ([]*domain.LearningProcess, error) {
	return s.learningRepo.GetByUserID(ctx, userID)
}

// GetLearningByID retrieves a learning process by ID
func (s *LearningService) GetLearningByID(ctx context.Context, id string) (*domain.LearningProcess, error) {
	return s.learningRepo.GetByID(ctx, id)
}

// CreateLearningFromRequest creates a learning process from topic and description
func (s *LearningService) CreateLearningFromRequest(ctx context.Context, userID, topic, description string) (*domain.LearningProcess, error) {
	// First, create a training request
	request := &domain.TrainingRequest{
		UserID:      userID,
		Topic:       topic,
		Description: description,
		Status:      domain.RequestApproved, // Auto-approve
	}

	if err := s.requestRepo.Create(ctx, request); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Get available mentor (with lowest workload)
	mentors, err := s.mentorRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get mentors: %w", err)
	}

	if len(mentors) == 0 {
		return nil, fmt.Errorf("no mentors available")
	}

	// Find mentor with lowest workload
	var selectedMentor *domain.Mentor
	minWorkload := 6 // Max is 5, so start with 6
	for _, mentor := range mentors {
		if mentor.Workload < minWorkload {
			minWorkload = mentor.Workload
			selectedMentor = mentor
		}
	}

	if selectedMentor == nil {
		return nil, fmt.Errorf("no suitable mentor found")
	}

	// Create learning process
	learning := &domain.LearningProcess{
		RequestID: request.ID,
		UserID:    userID,
		MentorID:  selectedMentor.ID,
		Status:    domain.LearningActive,
		StartDate: time.Now(),
		Plan:      []domain.LearningPlanItem{},
		Notes:     nil,
	}

	if err := s.learningRepo.Create(ctx, learning); err != nil {
		return nil, fmt.Errorf("failed to create learning process: %w", err)
	}

	// Update mentor workload
	if err := s.mentorRepo.UpdateWorkload(ctx, selectedMentor.ID, selectedMentor.Workload+1); err != nil {
		return nil, fmt.Errorf("failed to update mentor workload: %w", err)
	}

	// Reload to get full data with JOINs
	return s.learningRepo.GetByID(ctx, learning.ID)
}

// UpdateLearning updates full learning process (admin only)
func (s *LearningService) UpdateLearning(ctx context.Context, id string, topic string, description string, status domain.LearningStatus, plan []domain.LearningPlanItem, feedback *domain.Feedback, notes *string) (*domain.LearningProcess, error) {
	_, err := s.learningRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Create learning object for update
	learning := &domain.LearningProcess{
		Status:   status,
		Plan:     plan,
		Feedback: feedback,
		Notes:    notes,
	}

	// If status is completed and endDate is not set, it will be set in repository
	if status == domain.LearningCompleted && learning.EndDate == nil {
		now := time.Now()
		learning.EndDate = &now
	}

	if err := s.learningRepo.Update(ctx, id, learning); err != nil {
		return nil, fmt.Errorf("failed to update learning: %w", err)
	}

	// Note: topic and description are in the request, they don't change
	// They come from JOIN with training_requests table

	// Reload to get updated data with JOINs
	return s.learningRepo.GetByID(ctx, id)
}

// UpdatePlan updates learning plan
func (s *LearningService) UpdatePlan(ctx context.Context, id string, plan []domain.LearningPlanItem) (*domain.LearningProcess, error) {
	_, err := s.learningRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate plan items
	for _, item := range plan {
		if err := item.Validate(); err != nil {
			return nil, err
		}
	}

	if err := s.learningRepo.UpdatePlan(ctx, id, plan); err != nil {
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}

	// Reload to get updated data
	return s.learningRepo.GetByID(ctx, id)
}

// UpdateNotes updates learning notes
func (s *LearningService) UpdateNotes(ctx context.Context, id string, notes string) (*domain.LearningProcess, error) {
	_, err := s.learningRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.learningRepo.UpdateNotes(ctx, id, notes); err != nil {
		return nil, fmt.Errorf("failed to update notes: %w", err)
	}

	// Reload to get updated data
	return s.learningRepo.GetByID(ctx, id)
}

// CompleteLearning marks learning as completed with feedback
func (s *LearningService) CompleteLearning(ctx context.Context, id string, rating int, comment string) (*domain.LearningProcess, error) {
	learning, err := s.learningRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate
	if !learning.IsActive() {
		return nil, domain.ErrLearningNotActive
	}

	feedback := domain.Feedback{
		Rating:  rating,
		Comment: comment,
	}

	if err := feedback.Validate(); err != nil {
		return nil, err
	}

	if err := s.learningRepo.Complete(ctx, id, feedback); err != nil {
		return nil, fmt.Errorf("failed to complete learning: %w", err)
	}

	// Reload to get updated data
	return s.learningRepo.GetByID(ctx, id)
}

// CreateLearningProcess creates a new learning process for a request
func (s *LearningService) CreateLearningProcess(ctx context.Context, requestID, mentorID string) (*domain.LearningProcess, error) {
	// Get request
	request, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("request not found: %w", err)
	}

	// Get mentor
	mentor, err := s.mentorRepo.GetByID(ctx, mentorID)
	if err != nil {
		return nil, fmt.Errorf("mentor not found: %w", err)
	}

	// Create learning process
	learning := &domain.LearningProcess{
		RequestID: requestID,
		UserID:    request.UserID,
		MentorID:  mentorID,
		Status:    domain.LearningActive,
		StartDate: time.Now(),
		Plan:      []domain.LearningPlanItem{},
		Notes:     nil,
	}

	if err := s.learningRepo.Create(ctx, learning); err != nil {
		return nil, fmt.Errorf("failed to create learning process: %w", err)
	}

	// Update request status to approved
	if err := s.requestRepo.UpdateStatus(ctx, requestID, "approved"); err != nil {
		return nil, fmt.Errorf("failed to update request status: %w", err)
	}

	// Update mentor workload
	if err := s.mentorRepo.UpdateWorkload(ctx, mentorID, mentor.Workload+1); err != nil {
		return nil, fmt.Errorf("failed to update mentor workload: %w", err)
	}

	// Reload to get full data with JOINs
	return s.learningRepo.GetByID(ctx, learning.ID)
}
