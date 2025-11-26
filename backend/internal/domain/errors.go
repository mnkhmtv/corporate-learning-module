package domain

import "errors"

// Domain-level errors
var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden: insufficient permissions")

	// Mentor errors
	ErrMentorNotFound     = errors.New("mentor not found")
	ErrMentorNotAvailable = errors.New("mentor is not available (workload full)")

	// Training request errors
	ErrRequestNotFound        = errors.New("training request not found")
	ErrRequestAlreadyApproved = errors.New("request already approved")
	ErrRequestAlreadyRejected = errors.New("request already rejected")

	// Learning process errors
	ErrLearningNotFound      = errors.New("learning process not found")
	ErrLearningAlreadyExists = errors.New("learning process already exists for this request")
	ErrLearningNotActive     = errors.New("learning process is not active")
	ErrInvalidRating         = errors.New("rating must be between 1 and 5")
	ErrPlanItemNotFound      = errors.New("plan item not found")

	// Validation errors
	ErrInvalidInput    = errors.New("invalid input data")
	ErrEmptyField      = errors.New("required field is empty")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrWeakPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidWorkload = errors.New("workload must be between 0 and 5")
)
