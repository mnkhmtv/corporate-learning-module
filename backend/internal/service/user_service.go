package service

import (
	"context"
	"fmt"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by ID (admin only)
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetAllUsers retrieves all users (admin only)
func (s *UserService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.GetAll(ctx)
}

// UpdateUser updates user information (admin only)
func (s *UserService) UpdateUser(ctx context.Context, id string, name, email, department, jobTitle, telegram *string, password *string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != nil {
		user.Name = *name
	}
	if email != nil {
		user.Email = *email
	}
	if department != nil {
		user.Department = department
	}
	if jobTitle != nil {
		user.JobTitle = jobTitle
	}
	if telegram != nil {
		user.Telegram = telegram
	}

	// Update password if provided
	if password != nil && *password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// UpdateCurrentUser updates current user's profile
func (s *UserService) UpdateCurrentUser(ctx context.Context, userID string, name, email, department, jobTitle, telegram *string, password *string) (*domain.User, error) {
	// Same as UpdateUser but for current user
	return s.UpdateUser(ctx, userID, name, email, department, jobTitle, telegram, password)
}
