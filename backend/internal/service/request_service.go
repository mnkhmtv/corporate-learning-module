package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
)

type RequestService struct {
	requestRepo  domain.RequestRepository
	userRepo     domain.UserRepository
	mentorRepo   domain.MentorRepository
	learningRepo domain.LearningRepository
}

func NewRequestService(
	requestRepo domain.RequestRepository,
	userRepo domain.UserRepository,
	mentorRepo domain.MentorRepository,
	learningRepo domain.LearningRepository,
) *RequestService {
	return &RequestService{
		requestRepo:  requestRepo,
		userRepo:     userRepo,
		mentorRepo:   mentorRepo,
		learningRepo: learningRepo,
	}
}

// CreateRequest creates a new training request
func (s *RequestService) CreateRequest(ctx context.Context, userID, topic, description string) (*domain.TrainingRequest, error) {
	// Verify user exists
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	request := &domain.TrainingRequest{
		UserID:      userID,
		Topic:       topic,
		Description: description,
		Status:      domain.RequestPending,
	}

	if err := s.requestRepo.Create(ctx, request); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return request, nil
}

// GetAllRequests retrieves all training requests with optional status filter
func (s *RequestService) GetAllRequests(ctx context.Context, status *string) ([]*domain.TrainingRequest, error) {
	return s.requestRepo.GetAll(ctx, status)
}

// GetUserRequests retrieves all requests for a specific user
func (s *RequestService) GetUserRequests(ctx context.Context, userID string) ([]*domain.TrainingRequest, error) {
	return s.requestRepo.GetByUserID(ctx, userID)
}

// GetRequestByID retrieves a specific request by ID
func (s *RequestService) GetRequestByID(ctx context.Context, id string) (*domain.TrainingRequest, error) {
	return s.requestRepo.GetByID(ctx, id)
}

// UpdateRequest updates an existing training request
func (s *RequestService) UpdateRequest(ctx context.Context, id string, topic, description string) (*domain.TrainingRequest, error) {
	request, err := s.requestRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	request.Topic = topic
	request.Description = description

	if err := s.requestRepo.Update(ctx, request); err != nil {
		return nil, fmt.Errorf("failed to update request: %w", err)
	}

	return request, nil
}

// AssignMentor assigns a mentor to a request and creates learning process
func (s *RequestService) AssignMentor(ctx context.Context, requestID, mentorID string) (*domain.LearningProcess, error) {
	// Get request
	request, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	// Check if request is pending
	if request.Status != domain.RequestPending {
		return nil, fmt.Errorf("request is not pending")
	}

	// Get mentor
	mentor, err := s.mentorRepo.GetByID(ctx, mentorID)
	if err != nil {
		return nil, fmt.Errorf("mentor not found: %w", err)
	}

	// Check mentor workload
	if mentor.Workload >= 5 {
		return nil, fmt.Errorf("mentor has reached maximum workload")
	}

	// Approve request
	if err := s.requestRepo.UpdateStatus(ctx, requestID, string(domain.RequestApproved)); err != nil {
		return nil, fmt.Errorf("failed to approve request: %w", err)
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
		// Rollback request status
		_ = s.requestRepo.UpdateStatus(ctx, requestID, string(domain.RequestPending))
		return nil, fmt.Errorf("failed to create learning process: %w", err)
	}

	// Update mentor workload
	if err := s.mentorRepo.UpdateWorkload(ctx, mentorID, mentor.Workload+1); err != nil {
		return nil, fmt.Errorf("failed to update mentor workload: %w", err)
	}

	// Reload to get full data with JOINs
	return s.learningRepo.GetByID(ctx, learning.ID)
}
