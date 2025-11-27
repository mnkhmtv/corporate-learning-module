package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/domain"
	"github.com/mnkhmtv/corporate-learning-module/backend/internal/pkg/metrics"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RequestRepository struct {
	pool *pgxpool.Pool
}

func NewRequestRepository(pool *pgxpool.Pool) *RequestRepository {
	return &RequestRepository{pool: pool}
}

// Create inserts a new training request
func (r *RequestRepository) Create(ctx context.Context, req *domain.TrainingRequest) error {
	start := time.Now()

	query := `
		INSERT INTO training_requests (userId, topic, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, createdAt, updatedAt
	`

	err := r.pool.QueryRow(
		ctx, query,
		req.UserID, req.Topic, req.Description, req.Status,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)

	metrics.RecordDbQuery("requests.Create", time.Since(start), err)

	if err == nil {
		metrics.TrainingRequestsTotal.WithLabelValues(string(req.Status)).Inc()
	}

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	return nil
}

// GetByID retrieves a request by ID
func (r *RequestRepository) GetByID(ctx context.Context, id string) (*domain.TrainingRequest, error) {
	start := time.Now()

	query := `
		SELECT id, userId, topic, description, status, createdAt, updatedAt
		FROM training_requests
		WHERE id = $1
	`

	var req domain.TrainingRequest
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&req.ID, &req.UserID, &req.Topic, &req.Description,
		&req.Status, &req.CreatedAt, &req.UpdatedAt,
	)

	metrics.RecordDbQuery("requests.GetByID", time.Since(start), err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRequestNotFound
		}
		return nil, fmt.Errorf("failed to get request: %w", err)
	}

	return &req, nil
}

// GetByUserID retrieves all requests for a specific user
func (r *RequestRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.TrainingRequest, error) {
	start := time.Now()

	query := `
		SELECT id, userId, topic, description, status, createdAt, updatedAt
		FROM training_requests
		WHERE userId = $1
		ORDER BY createdAt DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)

	metrics.RecordDbQuery("requests.GetByUserID", time.Since(start), err)

	if err != nil {
		return nil, fmt.Errorf("failed to get requests by user: %w", err)
	}
	defer rows.Close()

	return r.scanRequests(rows)
}

// GetAll retrieves all requests with optional status filter
func (r *RequestRepository) GetAll(ctx context.Context, status *string) ([]*domain.TrainingRequest, error) {
	start := time.Now()

	var query string
	var args []interface{}

	if status != nil {
		query = `
			SELECT id, userId, topic, description, status, createdAt, updatedAt
			FROM training_requests
			WHERE status = $1
			ORDER BY createdAt DESC
		`
		args = append(args, *status)
	} else {
		query = `
			SELECT id, userId, topic, description, status, createdAt, updatedAt
			FROM training_requests
			ORDER BY createdAt DESC
		`
	}

	rows, err := r.pool.Query(ctx, query, args...)

	metrics.RecordDbQuery("requests.GetAll", time.Since(start), err)

	if err != nil {
		return nil, fmt.Errorf("failed to get all requests: %w", err)
	}
	defer rows.Close()

	return r.scanRequests(rows)
}

// Update updates an existing training request
func (r *RequestRepository) Update(ctx context.Context, req *domain.TrainingRequest) error {
	start := time.Now()

	query := `
		UPDATE training_requests
		SET topic = $2, description = $3
		WHERE id = $1
		RETURNING updatedAt
	`

	var updatedAt time.Time
	err := r.pool.QueryRow(ctx, query, req.ID, req.Topic, req.Description).Scan(&updatedAt)

	metrics.RecordDbQuery("requests.Update", time.Since(start), err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrRequestNotFound
		}
		return fmt.Errorf("failed to update request: %w", err)
	}

	req.UpdatedAt = updatedAt
	return nil
}

// UpdateStatus updates the status of a request
func (r *RequestRepository) UpdateStatus(ctx context.Context, id, status string) error {
	start := time.Now()

	query := `
		UPDATE training_requests
		SET status = $2
		WHERE id = $1
		RETURNING updatedAt
	`

	var updatedAt time.Time
	err := r.pool.QueryRow(ctx, query, id, status).Scan(&updatedAt)

	metrics.RecordDbQuery("requests.UpdateStatus", time.Since(start), err)

	if err == nil {
		metrics.TrainingRequestsTotal.WithLabelValues(status).Inc()
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrRequestNotFound
		}
		return fmt.Errorf("failed to update request status: %w", err)
	}

	return nil
}

// scanRequests is a helper to scan multiple rows into TrainingRequest slice
func (r *RequestRepository) scanRequests(rows pgx.Rows) ([]*domain.TrainingRequest, error) {

	var requests []*domain.TrainingRequest

	for rows.Next() {
		var req domain.TrainingRequest
		err := rows.Scan(
			&req.ID, &req.UserID, &req.Topic, &req.Description,
			&req.Status, &req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}
		requests = append(requests, &req)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return requests, nil
}
