package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MentorRepository struct {
	pool *pgxpool.Pool
}

func NewMentorRepository(pool *pgxpool.Pool) *MentorRepository {
	return &MentorRepository{pool: pool}
}

// Create inserts a new mentor
func (r *MentorRepository) Create(ctx context.Context, mentor *domain.Mentor) error {
	query := `
		INSERT INTO mentors (name, jobTitle, experience, workload, email, telegram)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, createdAt, updatedAt
	`

	err := r.pool.QueryRow(
		ctx, query,
		mentor.Name, mentor.JobTitle, mentor.Experience, mentor.Workload,
		mentor.Email, mentor.Telegram, mentor.Avatar,
	).Scan(&mentor.ID, &mentor.CreatedAt, &mentor.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create mentor: %w", err)
	}

	return nil
}

// GetByID retrieves a mentor by ID
func (r *MentorRepository) GetByID(ctx context.Context, id string) (*domain.Mentor, error) {
	query := `
		SELECT id, name, jobTitle, experience, workload, email, telegram, avatar, createdAt, updatedAt
		FROM mentors
		WHERE id = $1
	`

	var mentor domain.Mentor
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&mentor.ID, &mentor.Name, &mentor.JobTitle, &mentor.Experience,
		&mentor.Workload, &mentor.Email, &mentor.Telegram, &mentor.Avatar,
		&mentor.CreatedAt, &mentor.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrMentorNotFound
		}
		return nil, fmt.Errorf("failed to get mentor: %w", err)
	}

	return &mentor, nil
}

// GetAll retrieves all mentors with optional workload filter
func (r *MentorRepository) GetAll(ctx context.Context, maxWorkload *int) ([]*domain.Mentor, error) {
	var query string
	var args []interface{}

	if maxWorkload != nil {
		query = `
			SELECT id, name, job_title, experience, workload, email, telegram, avatar, createdAt, updatedAt
			FROM mentors
			WHERE workload <= $1
			ORDER BY workload ASC, name ASC
		`
		args = append(args, *maxWorkload)
	} else {
		query = `
			SELECT id, name, job_title, experience, workload, email, telegram, avatar, createdAt, updatedAt
			FROM mentors
			ORDER BY workload ASC, name ASC
		`
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get mentors: %w", err)
	}
	defer rows.Close()

	var mentors []*domain.Mentor
	for rows.Next() {
		var mentor domain.Mentor
		err := rows.Scan(
			&mentor.ID, &mentor.Name, &mentor.JobTitle, &mentor.Experience,
			&mentor.Workload, &mentor.Email, &mentor.Telegram, &mentor.Avatar,
			&mentor.CreatedAt, &mentor.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mentor: %w", err)
		}
		mentors = append(mentors, &mentor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return mentors, nil
}

// UpdateWorkload updates mentor's workload
func (r *MentorRepository) UpdateWorkload(ctx context.Context, id string, workload int) error {
	query := `
		UPDATE mentors
		SET workload = $2
		WHERE id = $1
		RETURNING updatedAt
	`

	var updatedAt string
	err := r.pool.QueryRow(ctx, query, id, workload).Scan(&updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrMentorNotFound
		}
		return fmt.Errorf("failed to update workload: %w", err)
	}

	return nil
}

// Delete removes a mentor
func (r *MentorRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM mentors WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete mentor: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrMentorNotFound
	}

	return nil
}
