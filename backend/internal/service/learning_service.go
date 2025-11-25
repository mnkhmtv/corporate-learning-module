package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
)

type LearningService struct {
	learningRepo domain.LearningRepository
	requestRepo  domain.RequestRepository
	mentorRepo   domain.MentorRepository
	mentorSvc    *MentorService
}

func NewLearningService(
	learningRepo domain.LearningRepository,
	requestRepo domain.RequestRepository,
	mentorRepo domain.MentorRepository,
	mentorSvc *MentorService,
) *LearningService {
	return &LearningService{
		learningRepo: learningRepo,
		requestRepo:  requestRepo,
		mentorRepo:   mentorRepo,
		mentorSvc:    mentorSvc,
	}
}

// AssignMentor creates a learning process by assigning a mentor to an approved request
func (s *LearningService) AssignMentor(ctx context.Context, requestID, mentorID string) (*domain.LearningProcess, error) {
	// Get and validate request
	request, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	if !request.IsApproved() {
		return nil, errors.New("request must be approved before assigning mentor")
	}

	// Get mentor
	mentor, err := s.mentorRepo.GetByID(ctx, mentorID)
	if err != nil {
		return nil, err
	}

	if !mentor.CanTakeStudent() {
		return nil, domain.ErrMentorNotAvailable
	}

	// Create learning process
	learning := &domain.LearningProcess{
		RequestID: requestID,
		UserID:    request.UserID,
		MentorID:  mentorID,
		Status:    domain.LearningActive,
		Plan:      []domain.LearningPlanItem{}, // Empty plan initially
	}

	if err := s.learningRepo.Create(ctx, learning); err != nil {
		return nil, fmt.Errorf("failed to create learning process: %w", err)
	}

	// Increment mentor workload
	if err := s.mentorSvc.IncrementMentorWorkload(ctx, mentorID); err != nil {
		return nil, fmt.Errorf("failed to update mentor workload: %w", err)
	}

	return learning, nil
}

// GetLearningByID retrieves a specific learning process
func (s *LearningService) GetLearningByID(ctx context.Context, learningID string) (*domain.LearningProcess, error) {
	return s.learningRepo.GetByID(ctx, learningID)
}

// GetUserLearnings retrieves all learning processes for a user
func (s *LearningService) GetUserLearnings(ctx context.Context, userID string) ([]*domain.LearningProcess, error) {
	return s.learningRepo.GetByUserID(ctx, userID)
}

// GetMentorLearnings retrieves all learning processes for a mentor
func (s *LearningService) GetMentorLearnings(ctx context.Context, mentorID string) ([]*domain.LearningProcess, error) {
	return s.learningRepo.GetByMentorID(ctx, mentorID)
}

// AddPlanItem adds a new item to the learning plan
func (s *LearningService) AddPlanItem(ctx context.Context, learningID, text string) (*domain.LearningProcess, error) {
	learning, err := s.learningRepo.GetByID(ctx, learningID)
	if err != nil {
		return nil, err
	}

	if !learning.IsActive() {
		return nil, domain.ErrLearningNotActive
	}

	// Create new plan item
	item, err := domain.NewLearningPlanItem(text)
	if err != nil {
		return nil, err
	}

	// Add to plan
	if err := learning.AddPlanItem(*item); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.learningRepo.UpdatePlan(ctx, learning.ID, learning.Plan); err != nil {
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}

	return learning, nil
}

// UpdatePlanItem updates an existing plan item
func (s *LearningService) UpdatePlanItem(ctx context.Context, learningID, itemID, text string, completed bool) error {
	learning, err := s.learningRepo.GetByID(ctx, learningID)
	if err != nil {
		return err
	}

	if err := learning.UpdatePlanItem(itemID, text, completed); err != nil {
		return err
	}

	return s.learningRepo.UpdatePlan(ctx, learning.ID, learning.Plan)
}

// TogglePlanItem toggles completion status of a plan item
func (s *LearningService) TogglePlanItem(ctx context.Context, learningID, itemID string) error {
	learning, err := s.learningRepo.GetByID(ctx, learningID)
	if err != nil {
		return err
	}

	if err := learning.TogglePlanItem(itemID); err != nil {
		return err
	}

	return s.learningRepo.UpdatePlan(ctx, learning.ID, learning.Plan)
}

// RemovePlanItem removes an item from the plan
func (s *LearningService) RemovePlanItem(ctx context.Context, learningID, itemID string) error {
	learning, err := s.learningRepo.GetByID(ctx, learningID)
	if err != nil {
		return err
	}

	if err := learning.RemovePlanItem(itemID); err != nil {
		return err
	}

	return s.learningRepo.UpdatePlan(ctx, learning.ID, learning.Plan)
}

// CompleteLearning marks the learning process as completed with feedback
func (s *LearningService) CompleteLearning(ctx context.Context, learningID string, rating int, comment string) error {
	learning, err := s.learningRepo.GetByID(ctx, learningID)
	if err != nil {
		return err
	}

	// Validate and complete
	if err := learning.Complete(rating, comment); err != nil {
		return err
	}

	// Update in database
	if err := s.learningRepo.Complete(ctx, learning.ID, rating, comment); err != nil {
		return fmt.Errorf("failed to complete learning: %w", err)
	}

	// Decrement mentor workload
	if err := s.mentorSvc.DecrementMentorWorkload(ctx, learning.MentorID); err != nil {
		return fmt.Errorf("failed to update mentor workload: %w", err)
	}

	return nil
}

// GetLearningProgress calculates progress percentage
func (s *LearningService) GetLearningProgress(ctx context.Context, learningID string) (float64, error) {
	learning, err := s.learningRepo.GetByID(ctx, learningID)
	if err != nil {
		return 0, err
	}

	return learning.GetProgress(), nil
}

// UpdatePlan updates the entire learning plan
func (s *LearningService) UpdatePlan(ctx context.Context, learningID string, planItems []domain.LearningPlanItem) (*domain.LearningProcess, error) {
	learning, err := s.learningRepo.GetByID(ctx, learningID)
	if err != nil {
		return nil, err
	}

	if !learning.IsActive() {
		return nil, domain.ErrLearningNotActive
	}

	// Validate plan elements
	for _, item := range planItems {
		if err := item.Validate(); err != nil {
			return nil, err
		}
	}

	// Update plan
	if err := s.learningRepo.UpdatePlan(ctx, learning.ID, planItems); err != nil {
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}

	// Receive updated learning process
	return s.learningRepo.GetByID(ctx, learningID)
}

// UpdateNotes updates notes for a learning process
func (s *LearningService) UpdateNotes(ctx context.Context, learningID, notes string) (*domain.LearningProcess, error) {
	learning, err := s.learningRepo.GetByID(ctx, learningID)
	if err != nil {
		return nil, err
	}

	if !learning.IsActive() {
		return nil, domain.ErrLearningNotActive
	}

	// Обновить заметки
	if err := s.learningRepo.UpdateNotes(ctx, learning.ID, notes); err != nil {
		return nil, fmt.Errorf("failed to update notes: %w", err)
	}

	// Получить обновленный learning process
	return s.learningRepo.GetByID(ctx, learningID)
}
