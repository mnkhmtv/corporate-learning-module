package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LearningRepository struct {
	pool *pgxpool.Pool
}

func NewLearningRepository(pool *pgxpool.Pool) *LearningRepository {
	return &LearningRepository{pool: pool}
}

// Create inserts a new learning process
func (r *LearningRepository) Create(ctx context.Context, learning *domain.LearningProcess) error {
	planJSON, err := json.Marshal(learning.Plan)
	if err != nil {
		return fmt.Errorf("failed to marshal plan: %w", err)
	}

	query := `
		INSERT INTO learning_processes 
		(requestId, userId, mentorId, status, plan, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, createdAt, updatedAt
	`

	err = r.pool.QueryRow(
		ctx, query,
		learning.RequestID, learning.UserID, learning.MentorID,
		learning.MentorName, learning.MentorEmail, learning.Status,
		planJSON, learning.Notes,
	).Scan(&learning.ID, &learning.CreatedAt, &learning.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create learning process: %w", err)
	}

	return nil
}

// GetByID retrieves a learning process by ID
func (r *LearningRepository) GetByID(ctx context.Context, id string) (*domain.LearningProcess, error) {
	query := `
		SELECT id, requestId, userId, mentorId, 
		       status, plan, notes, feedbackRating, feedbackComment,
		       createdAt, updatedAt, completedAt
		FROM learning_processes
		WHERE id = $1
	`

	var learning domain.LearningProcess
	var planJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&learning.ID, &learning.RequestID, &learning.UserID,
		&learning.MentorID, &learning.MentorName, &learning.MentorEmail,
		&learning.Status, &planJSON, &learning.Notes,
		&learning.FeedbackRating, &learning.FeedbackComment,
		&learning.CreatedAt, &learning.UpdatedAt, &learning.CompletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLearningNotFound
		}
		return nil, fmt.Errorf("failed to get learning process: %w", err)
	}

	if err := json.Unmarshal(planJSON, &learning.Plan); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plan: %w", err)
	}

	return &learning, nil
}

// GetByUserID retrieves all learning processes for a user
func (r *LearningRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.LearningProcess, error) {
	query := `
		SELECT id, requestId, userId, mentorId, 
		       status, plan, notes, feedbackRating, feedbackComment,
		       createdAt, updatedAt, completedAt
		FROM learning_processes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get learning processes by user: %w", err)
	}
	defer rows.Close()

	return r.scanLearningProcesses(rows)
}

// GetByMentorID retrieves all learning processes for a mentor
func (r *LearningRepository) GetByMentorID(ctx context.Context, mentorID string) ([]*domain.LearningProcess, error) {
	query := `
		SELECT id, requestId, userId, mentorId, 
		       status, plan, notes, feedbackRating, feedbackComment,
		       createdAt, updatedAt, completedAt
		FROM learning_processes
		WHERE mentor_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, mentorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get learning processes by mentor: %w", err)
	}
	defer rows.Close()

	return r.scanLearningProcesses(rows)
}

// UpdatePlan updates the learning plan
func (r *LearningRepository) UpdatePlan(ctx context.Context, id string, plan []domain.LearningPlanItem) error {
	planJSON, err := json.Marshal(plan)
	if err != nil {
		return fmt.Errorf("failed to marshal plan: %w", err)
	}

	query := `
		UPDATE learning_processes
		SET plan = $2
		WHERE id = $1
		RETURNING updatedAt
	`

	var updatedAt string
	err = r.pool.QueryRow(ctx, query, id, planJSON).Scan(&updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrLearningNotFound
		}
		return fmt.Errorf("failed to update plan: %w", err)
	}

	return nil
}

// Complete marks learning as completed with feedback
func (r *LearningRepository) Complete(ctx context.Context, id string, rating int, comment string) error {
	query := `
		UPDATE learning_processes
		SET status = 'completed',
		    feedbackRating = $2,
		    feedbackComment = $3,
		    completedAt = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updatedAt
	`

	var updatedAt string
	err := r.pool.QueryRow(ctx, query, id, rating, comment).Scan(&updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrLearningNotFound
		}
		return fmt.Errorf("failed to complete learning: %w", err)
	}

	return nil
}

// scanLearningProcesses is a helper to scan multiple rows
func (r *LearningRepository) scanLearningProcesses(rows pgx.Rows) ([]*domain.LearningProcess, error) {
	var learnings []*domain.LearningProcess

	for rows.Next() {
		var learning domain.LearningProcess
		var planJSON []byte

		err := rows.Scan(
			&learning.ID, &learning.RequestID, &learning.UserID,
			&learning.MentorID, &learning.MentorName, &learning.MentorEmail,
			&learning.Status, &planJSON, &learning.Notes,
			&learning.FeedbackRating, &learning.FeedbackComment,
			&learning.CreatedAt, &learning.UpdatedAt, &learning.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan learning process: %w", err)
		}

		if err := json.Unmarshal(planJSON, &learning.Plan); err != nil {
			return nil, fmt.Errorf("failed to unmarshal plan: %w", err)
		}

		learnings = append(learnings, &learning)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return learnings, nil
}
