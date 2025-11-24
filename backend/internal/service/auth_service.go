package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  domain.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthService(userRepo domain.UserRepository, jwtSecret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, name, email, password, role, department, jobTitle, telegram string) (*domain.User, error) {
	// Validate input
	if err := s.validateRegistration(name, email, password); err != nil {
		return nil, err
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Convert role string to domain.UserRole
	userRole := domain.UserRole(role)
	if userRole != domain.RoleEmployee && userRole != domain.RoleAdmin {
		userRole = domain.RoleEmployee // Default to employee
	}

	// Create user
	user := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         userRole,
		Department:   stringToPtr(department),
		JobTitle:     stringToPtr(jobTitle),
		Telegram:     stringToPtr(telegram),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns JWT token
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", nil, domain.ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, domain.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, string(user.Role), s.jwtSecret, s.tokenTTL)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

// GetUserByID retrieves user information
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// validateRegistration validates registration input
func (s *AuthService) validateRegistration(name, email, password string) error {
	if name == "" {
		return fmt.Errorf("%w: name is required", domain.ErrInvalidInput)
	}
	if email == "" {
		return fmt.Errorf("%w: email is required", domain.ErrInvalidInput)
	}
	if len(password) < 8 {
		return domain.ErrWeakPassword
	}
	// Add email format validation if needed
	return nil
}

// hashPassword generates bcrypt hash from plain password
func (s *AuthService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Helper function to convert string to pointer
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
