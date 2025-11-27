package domain

import "context"

// UserRepository defines methods for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// RequestRepository defines methods for training request data access
type RequestRepository interface {
	Create(ctx context.Context, request *TrainingRequest) error
	GetByID(ctx context.Context, id string) (*TrainingRequest, error)
	GetByUserID(ctx context.Context, userID string) ([]*TrainingRequest, error)
	GetAll(ctx context.Context, status *string) ([]*TrainingRequest, error)
	Update(ctx context.Context, request *TrainingRequest) error
	UpdateStatus(ctx context.Context, id, status string) error
}

// MentorRepository defines methods for mentor data access
type MentorRepository interface {
	Create(ctx context.Context, mentor *Mentor) error
	GetByID(ctx context.Context, id string) (*Mentor, error)
	GetAll(ctx context.Context, maxWorkload *int) ([]*Mentor, error)
	Update(ctx context.Context, mentor *Mentor) error
	UpdateWorkload(ctx context.Context, id string, workload int) error
	Delete(ctx context.Context, id string) error
}

// LearningRepository defines methods for learning process data access
type LearningRepository interface {
	Create(ctx context.Context, learning *LearningProcess) error
	GetByID(ctx context.Context, id string) (*LearningProcess, error)
	GetByUserID(ctx context.Context, userID string) ([]*LearningProcess, error)
	GetByMentorID(ctx context.Context, mentorID string) ([]*LearningProcess, error)
	GetAll(ctx context.Context) ([]*LearningProcess, error)
	UpdatePlan(ctx context.Context, id string, plan []LearningPlanItem) error
	UpdateNotes(ctx context.Context, id string, notes string) error
	Update(ctx context.Context, id string, learning *LearningProcess) error
	Complete(ctx context.Context, id string, feedback Feedback) error
}
