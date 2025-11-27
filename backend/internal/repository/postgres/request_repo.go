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
func (r *RequestRepository) Create(ctx context.Context, request *domain.TrainingRequest) error {
	start := time.Now()

	query := `
		INSERT INTO training_requests (userId, topic, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, createdAt, updatedAt
	`

	err := r.pool.QueryRow(
		ctx, query,
		request.UserID, request.Topic, request.Description, request.Status,
	).Scan(&request.ID, &request.CreatedAt, &request.UpdatedAt)

	metrics.RecordDbQuery("requests.Create", time.Since(start), err)

	if err == nil {
		metrics.TrainingRequestsTotal.WithLabelValues(string(request.Status)).Inc()
	}

	if err != nil {
		return fmt.Errorf("failed to create training request: %w", err)
	}

	return nil
}

// GetByID retrieves a training request by ID with user data
func (r *RequestRepository) GetByID(ctx context.Context, id string) (*domain.TrainingRequest, error) {
	start := time.Now()

	query := `
		SELECT 
			r.id, r.userId, r.topic, r.description, r.status, r.createdAt, r.updatedAt,
			u.name AS userName,
			u.jobTitle AS userJobTitle,
			u.telegram AS userTelegram
		FROM training_requests r
		INNER JOIN users u ON r.userId = u.id
		WHERE r.id = $1
	`

	var request domain.TrainingRequest
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&request.ID, &request.UserID, &request.Topic, &request.Description,
		&request.Status, &request.CreatedAt, &request.UpdatedAt,
		&request.UserName, &request.UserJobTitle, &request.UserTelegram,
	)

	metrics.RecordDbQuery("requests.GetByID", time.Since(start), err)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRequestNotFound
		}
		return nil, fmt.Errorf("failed to get training request: %w", err)
	}

	return &request, nil
}

// GetByUserID retrieves all training requests for a user
func (r *RequestRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.TrainingRequest, error) {
	start := time.Now()

	query := `
		SELECT 
			r.id, r.userId, r.topic, r.description, r.status, r.createdAt, r.updatedAt,
			u.name AS userName,
			u.jobTitle AS userJobTitle,
			u.telegram AS userTelegram
		FROM training_requests r
		INNER JOIN users u ON r.userId = u.id
		WHERE r.userId = $1
		ORDER BY r.createdAt DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)

	metrics.RecordDbQuery("requests.GetByUserID", time.Since(start), err)

	if err != nil {
		return nil, fmt.Errorf("failed to get user training requests: %w", err)
	}
	defer rows.Close()

	return r.scanRequests(rows)
}

// GetAll retrieves all training requests with optional status filter
func (r *RequestRepository) GetAll(ctx context.Context, status *string) ([]*domain.TrainingRequest, error) {
	start := time.Now()

	query := `
		SELECT 
			r.id, r.userId, r.topic, r.description, r.status, r.createdAt, r.updatedAt,
			u.name AS userName,
			u.jobTitle AS userJobTitle,
			u.telegram AS userTelegram
		FROM training_requests r
		INNER JOIN users u ON r.userId = u.id
	`

	var rows pgx.Rows
	var err error

	if status != nil {
		query += " WHERE r.status = $1 ORDER BY r.createdAt DESC"
		rows, err = r.pool.Query(ctx, query, *status)
	} else {
		query += " ORDER BY r.createdAt DESC"
		rows, err = r.pool.Query(ctx, query)
	}

	metrics.RecordDbQuery("requests.GetAll", time.Since(start), err)

	if err != nil {
		return nil, fmt.Errorf("failed to get all training requests: %w", err)
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

// UpdateStatus updates the status of a training request
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

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrRequestNotFound
		}
		return fmt.Errorf("failed to update request status: %w", err)
	}

	return nil
}

// scanRequests is a helper to scan multiple rows
func (r *RequestRepository) scanRequests(rows pgx.Rows) ([]*domain.TrainingRequest, error) {
	var requests []*domain.TrainingRequest

	for rows.Next() {
		var request domain.TrainingRequest
		err := rows.Scan(
			&request.ID, &request.UserID, &request.Topic, &request.Description,
			&request.Status, &request.CreatedAt, &request.UpdatedAt,
			&request.UserName, &request.UserJobTitle, &request.UserTelegram,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan training request: %w", err)
		}
		requests = append(requests, &request)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return requests, nil
}
