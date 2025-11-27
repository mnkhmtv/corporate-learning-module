package service

import (
	"context"
	"errors"
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
	// Validate input
	if topic == "" {
		return nil, fmt.Errorf("%w: topic is required", domain.ErrInvalidInput)
	}
	if description == "" {
		return nil, fmt.Errorf("%w: description is required", domain.ErrInvalidInput)
	}

	// Verify user exists
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return nil, err
	}

	// Create request
	request := &domain.TrainingRequest{
		UserID:      userID,
		Topic:       topic,
		Description: description,
		Status:      domain.StatusPending,
	}

	if err := s.requestRepo.Create(ctx, request); err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return request, nil
}

// GetRequestByID retrieves a specific training request
func (s *RequestService) GetRequestByID(ctx context.Context, requestID string) (*domain.TrainingRequest, error) {
	return s.requestRepo.GetByID(ctx, requestID)
}

// GetUserRequests retrieves all requests for a specific user
func (s *RequestService) GetUserRequests(ctx context.Context, userID string) ([]*domain.TrainingRequest, error) {
	return s.requestRepo.GetByUserID(ctx, userID)
}

// GetAllRequests retrieves all requests (admin only)
func (s *RequestService) GetAllRequests(ctx context.Context, status *string) ([]*domain.TrainingRequest, error) {
	return s.requestRepo.GetAll(ctx, status)
}

// ApproveRequest approves a training request
func (s *RequestService) ApproveRequest(ctx context.Context, requestID string) error {
	request, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	if request.IsApproved() {
		return domain.ErrRequestAlreadyApproved
	}

	if request.IsRejected() {
		return errors.New("cannot approve a rejected request")
	}

	request.Approve()
	return s.requestRepo.UpdateStatus(ctx, request.ID, string(request.Status))
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

// RejectRequest rejects a training request
func (s *RequestService) RejectRequest(ctx context.Context, requestID string) error {
	request, err := s.requestRepo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}

	if request.IsRejected() {
		return domain.ErrRequestAlreadyRejected
	}

	if request.IsApproved() {
		return errors.New("cannot reject an approved request")
	}

	request.Reject()
	return s.requestRepo.UpdateStatus(ctx, request.ID, string(request.Status))
}
