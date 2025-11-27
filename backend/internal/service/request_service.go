package service

import (
	"context"
	"fmt"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
)

type RequestService struct {
	requestRepo domain.RequestRepository
	userRepo    domain.UserRepository
}

func NewRequestService(requestRepo domain.RequestRepository, userRepo domain.UserRepository) *RequestService {
	return &RequestService{
		requestRepo: requestRepo,
		userRepo:    userRepo,
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
