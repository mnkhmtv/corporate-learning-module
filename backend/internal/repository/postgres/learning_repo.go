package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/pkg/metrics"

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
	start := time.Now()

	planJSON, err := json.Marshal(learning.Plan)
	if err != nil {
		return fmt.Errorf("failed to marshal plan: %w", err)
	}

	query := `
		INSERT INTO learning_processes 
		(requestId, userId, mentorId, status, startDate, plan, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, startDate, createdAt, updatedAt
	`

	err = r.pool.QueryRow(
		ctx, query,
		learning.RequestID, learning.UserID, learning.MentorID,
		learning.Status, learning.StartDate, planJSON, learning.Notes,
	).Scan(&learning.ID, &learning.StartDate, &learning.CreatedAt, &learning.UpdatedAt)

	metrics.RecordDbQuery("learning.Create", time.Since(start), err)

	if err == nil {
		metrics.LearningProcessesActive.Inc()
	}

	if err != nil {
		return fmt.Errorf("failed to create learning process: %w", err)
	}

	return nil
}

// GetByID retrieves a learning process by ID with JOINs
func (r *LearningRepository) GetByID(ctx context.Context, id string) (*domain.LearningProcess, error) {
	start := time.Now()

	query := `
		SELECT 
			lp.id, lp.requestId, lp.userId, lp.mentorId,
			lp.status, lp.startDate, lp.endDate,
			lp.plan, lp.feedback, lp.notes,
			lp.createdAt, lp.updatedAt,
			r.topic AS requestTopic,
			r.description AS requestDescription,
			u.name AS userName,
			m.name AS mentorName,
			m.telegram AS mentorTelegram,
			m.jobTitle AS mentorJobTitle,
			m.experience AS mentorExperience
		FROM learning_processes lp
		INNER JOIN training_requests r ON lp.requestId = r.id
		INNER JOIN users u ON lp.userId = u.id
		INNER JOIN mentors m ON lp.mentorId = m.id
		WHERE lp.id = $1
	`

	var learning domain.LearningProcess
	var planJSON []byte
	var feedbackJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&learning.ID, &learning.RequestID, &learning.UserID, &learning.MentorID,
		&learning.Status, &learning.StartDate, &learning.EndDate,
		&planJSON, &feedbackJSON, &learning.Notes,
		&learning.CreatedAt, &learning.UpdatedAt,
		&learning.RequestTopic, &learning.RequestDescription,
		&learning.UserName,
		&learning.MentorName, &learning.MentorTelegram,
		&learning.MentorJobTitle, &learning.MentorExperience,
	)

	metrics.RecordDbQuery("learning.GetByID", time.Since(start), err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrLearningNotFound
		}
		return nil, fmt.Errorf("failed to get learning process: %w", err)
	}

	// Unmarshal plan
	if err := json.Unmarshal(planJSON, &learning.Plan); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plan: %w", err)
	}

	// Unmarshal feedback (может быть NULL)
	if feedbackJSON != nil {
		var feedback domain.Feedback
		if err := json.Unmarshal(feedbackJSON, &feedback); err != nil {
			return nil, fmt.Errorf("failed to unmarshal feedback: %w", err)
		}
		learning.Feedback = &feedback
	}

	return &learning, nil
}

// GetByUserID retrieves all learning processes for a user
func (r *LearningRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.LearningProcess, error) {
	start := time.Now()

	query := `
		SELECT 
			lp.id, lp.requestId, lp.userId, lp.mentorId,
			lp.status, lp.startDate, lp.endDate,
			lp.plan, lp.feedback, lp.notes,
			lp.createdAt, lp.updatedAt,
			r.topic AS requestTopic,
			r.description AS requestDescription,
			u.name AS userName,
			m.name AS mentorName,
			m.telegram AS mentorTelegram,
			m.jobTitle AS mentorJobTitle,
			m.experience AS mentorExperience
		FROM learning_processes lp
		INNER JOIN training_requests r ON lp.requestId = r.id
		INNER JOIN users u ON lp.userId = u.id
		INNER JOIN mentors m ON lp.mentorId = m.id
		WHERE lp.userId = $1
		ORDER BY lp.createdAt DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)

	metrics.RecordDbQuery("learning.GetByUserID", time.Since(start), err)

	if err != nil {
		return nil, fmt.Errorf("failed to get learning processes by user: %w", err)
	}
	defer rows.Close()

	return r.scanLearningProcesses(rows)
}

// GetByMentorID retrieves all learning processes for a mentor
func (r *LearningRepository) GetByMentorID(ctx context.Context, mentorID string) ([]*domain.LearningProcess, error) {
	start := time.Now()

	query := `
		SELECT 
			lp.id, lp.requestId, lp.userId, lp.mentorId,
			lp.status, lp.startDate, lp.endDate,
			lp.plan, lp.feedback, lp.notes,
			lp.createdAt, lp.updatedAt,
			r.topic AS requestTopic,
			r.description AS requestDescription,
			u.name AS userName,
			m.name AS mentorName,
			m.telegram AS mentorTelegram,
			m.jobTitle AS mentorJobTitle,
			m.experience AS mentorExperience
		FROM learning_processes lp
		INNER JOIN training_requests r ON lp.requestId = r.id
		INNER JOIN users u ON lp.userId = u.id
		INNER JOIN mentors m ON lp.mentorId = m.id
		WHERE lp.mentorId = $1
		ORDER BY lp.createdAt DESC
	`

	rows, err := r.pool.Query(ctx, query, mentorID)

	metrics.RecordDbQuery("learning.GetByMentorID", time.Since(start), err)

	if err != nil {
		return nil, fmt.Errorf("failed to get learning processes by mentor: %w", err)
	}
	defer rows.Close()

	return r.scanLearningProcesses(rows)
}

// GetAll retrieves all learning processes (admin only)
func (r *LearningRepository) GetAll(ctx context.Context) ([]*domain.LearningProcess, error) {
	start := time.Now()

	query := `
		SELECT 
			lp.id, lp.requestId, lp.userId, lp.mentorId,
			lp.status, lp.startDate, lp.endDate,
			lp.plan, lp.feedback, lp.notes,
			lp.createdAt, lp.updatedAt,
			r.topic AS requestTopic,
			r.description AS requestDescription,
			u.name AS userName,
			m.name AS mentorName,
			m.telegram AS mentorTelegram,
			m.jobTitle AS mentorJobTitle,
			m.experience AS mentorExperience
		FROM learning_processes lp
		INNER JOIN training_requests r ON lp.requestId = r.id
		INNER JOIN users u ON lp.userId = u.id
		INNER JOIN mentors m ON lp.mentorId = m.id
		ORDER BY lp.createdAt DESC
	`

	rows, err := r.pool.Query(ctx, query)

	metrics.RecordDbQuery("learning.GetAll", time.Since(start), err)

	if err != nil {
		return nil, fmt.Errorf("failed to get all learning processes: %w", err)
	}
	defer rows.Close()

	return r.scanLearningProcesses(rows)
}

// UpdateMentor updates the mentor for a learning process
func (r *LearningRepository) UpdateMentor(ctx context.Context, learningID, mentorID string) error {
	start := time.Now()

	query := `
		UPDATE learning_processes
		SET mentorId = $2
		WHERE id = $1
		RETURNING updatedAt
	`

	var updatedAt time.Time
	err := r.pool.QueryRow(ctx, query, learningID, mentorID).Scan(&updatedAt)

	metrics.RecordDbQuery("learning.UpdateMentor", time.Since(start), err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrLearningNotFound
		}
		return fmt.Errorf("failed to update mentor: %w", err)
	}

	return nil
}

// UpdatePlan updates the learning plan
func (r *LearningRepository) UpdatePlan(ctx context.Context, id string, plan []domain.LearningPlanItem) error {
	start := time.Now()

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

	var updatedAt time.Time
	err = r.pool.QueryRow(ctx, query, id, planJSON).Scan(&updatedAt)

	metrics.RecordDbQuery("learning.UpdatePlan", time.Since(start), err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrLearningNotFound
		}
		return fmt.Errorf("failed to update plan: %w", err)
	}

	return nil
}

// UpdateNotes updates notes for a learning process
func (r *LearningRepository) UpdateNotes(ctx context.Context, id string, notes string) error {
	start := time.Now()

	query := `
		UPDATE learning_processes
		SET notes = $2
		WHERE id = $1
		RETURNING updatedAt
	`

	var updatedAt time.Time
	var notesPtr *string
	if notes == "" {
		notesPtr = nil
	} else {
		notesPtr = &notes
	}

	err := r.pool.QueryRow(ctx, query, id, notesPtr).Scan(&updatedAt)

	metrics.RecordDbQuery("learning.UpdateNotes", time.Since(start), err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrLearningNotFound
		}
		return fmt.Errorf("failed to update notes: %w", err)
	}

	return nil
}

// Update updates full learning process info (admin only)
func (r *LearningRepository) Update(ctx context.Context, id string, learning *domain.LearningProcess) error {
	start := time.Now()

	planJSON, err := json.Marshal(learning.Plan)
	if err != nil {
		return fmt.Errorf("failed to marshal plan: %w", err)
	}

	var feedbackJSON []byte
	if learning.Feedback != nil {
		feedbackJSON, err = json.Marshal(learning.Feedback)
		if err != nil {
			return fmt.Errorf("failed to marshal feedback: %w", err)
		}
	}

	query := `
		UPDATE learning_processes
		SET status = $2,
		    plan = $3,
		    feedback = $4,
		    notes = $5,
		    endDate = $6
		WHERE id = $1
		RETURNING updatedAt
	`

	var updatedAt time.Time
	err = r.pool.QueryRow(ctx, query, id, learning.Status, planJSON, feedbackJSON, learning.Notes, learning.EndDate).Scan(&updatedAt)

	metrics.RecordDbQuery("learning.Update", time.Since(start), err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrLearningNotFound
		}
		return fmt.Errorf("failed to update learning: %w", err)
	}

	return nil
}

// Complete marks learning as completed with feedback
func (r *LearningRepository) Complete(ctx context.Context, id string, feedback domain.Feedback) error {
	start := time.Now()

	feedbackJSON, err := json.Marshal(feedback)
	if err != nil {
		return fmt.Errorf("failed to marshal feedback: %w", err)
	}

	query := `
		UPDATE learning_processes
		SET status = 'completed',
		    feedback = $2,
		    endDate = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING endDate, updatedAt
	`

	var endDate, updatedAt time.Time
	err = r.pool.QueryRow(ctx, query, id, feedbackJSON).Scan(&endDate, &updatedAt)

	metrics.RecordDbQuery("learning.Complete", time.Since(start), err)

	if err == nil {
		metrics.LearningProcessesActive.Dec()
		metrics.LearningProcessesCompleted.Inc()
		metrics.FeedbackRatingSum.Add(float64(feedback.Rating))
		metrics.FeedbackRatingCount.Inc()
	}

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
		var feedbackJSON []byte

		err := rows.Scan(
			&learning.ID, &learning.RequestID, &learning.UserID, &learning.MentorID,
			&learning.Status, &learning.StartDate, &learning.EndDate,
			&planJSON, &feedbackJSON, &learning.Notes,
			&learning.CreatedAt, &learning.UpdatedAt,
			&learning.RequestTopic, &learning.RequestDescription,
			&learning.UserName,
			&learning.MentorName, &learning.MentorTelegram,
			&learning.MentorJobTitle, &learning.MentorExperience,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan learning process: %w", err)
		}

		if err := json.Unmarshal(planJSON, &learning.Plan); err != nil {
			return nil, fmt.Errorf("failed to unmarshal plan: %w", err)
		}

		if feedbackJSON != nil {
			var feedback domain.Feedback
			if err := json.Unmarshal(feedbackJSON, &feedback); err != nil {
				return nil, fmt.Errorf("failed to unmarshal feedback: %w", err)
			}
			learning.Feedback = &feedback
		}

		learnings = append(learnings, &learning)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return learnings, nil
}
