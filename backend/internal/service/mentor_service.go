package service

import (
	"context"
	"fmt"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
)

type MentorService struct {
	mentorRepo domain.MentorRepository
}

func NewMentorService(mentorRepo domain.MentorRepository) *MentorService {
	return &MentorService{
		mentorRepo: mentorRepo,
	}
}

// CreateMentor creates a new mentor
func (s *MentorService) CreateMentor(ctx context.Context, name, jobTitle, experience, email, telegram string) (*domain.Mentor, error) {
	// Validate input
	if name == "" || jobTitle == "" || email == "" {
		return nil, fmt.Errorf("%w: name, jobTitle, and email are required", domain.ErrInvalidInput)
	}

	mentor := &domain.Mentor{
		Name:       name,
		JobTitle:   jobTitle,
		Experience: stringToPtr(experience),
		Workload:   0,
		Email:      email,
		Telegram:   stringToPtr(telegram),
	}

	if err := s.mentorRepo.Create(ctx, mentor); err != nil {
		return nil, fmt.Errorf("failed to create mentor: %w", err)
	}

	return mentor, nil
}

// GetMentorByID retrieves a specific mentor
func (s *MentorService) GetMentorByID(ctx context.Context, mentorID string) (*domain.Mentor, error) {
	return s.mentorRepo.GetByID(ctx, mentorID)
}

// GetAvailableMentors retrieves mentors with workload less than maximum
func (s *MentorService) GetAvailableMentors(ctx context.Context) ([]*domain.Mentor, error) {
	maxWorkload := 4 // Only mentors with workload <= 4 can accept new students
	return s.mentorRepo.GetAll(ctx, &maxWorkload)
}

// GetAllMentors retrieves all mentors
func (s *MentorService) GetAllMentors(ctx context.Context) ([]*domain.Mentor, error) {
	return s.mentorRepo.GetAll(ctx, nil)
}

// IncrementMentorWorkload increases a mentor's workload
func (s *MentorService) IncrementMentorWorkload(ctx context.Context, mentorID string) error {
	mentor, err := s.mentorRepo.GetByID(ctx, mentorID)
	if err != nil {
		return err
	}

	if !mentor.CanTakeStudent() {
		return domain.ErrMentorNotAvailable
	}

	mentor.IncrementWorkload()
	return s.mentorRepo.UpdateWorkload(ctx, mentor.ID, mentor.Workload)
}

// DecrementMentorWorkload decreases a mentor's workload
func (s *MentorService) DecrementMentorWorkload(ctx context.Context, mentorID string) error {
	mentor, err := s.mentorRepo.GetByID(ctx, mentorID)
	if err != nil {
		return err
	}

	mentor.DecrementWorkload()
	return s.mentorRepo.UpdateWorkload(ctx, mentor.ID, mentor.Workload)
}
