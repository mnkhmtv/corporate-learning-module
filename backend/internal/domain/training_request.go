package domain

import "time"

// RequestStatus defines the lifecycle status of a training request
type RequestStatus string

const (
	StatusPending  RequestStatus = "pending"
	StatusApproved RequestStatus = "approved"
	StatusRejected RequestStatus = "rejected"
)

// TrainingRequest represents an employee's request to learn a topic
type TrainingRequest struct {
	ID          uint64        `json:"id"`
	UserID      uint64        `json:"userId"`
	Topic       string        `json:"topic"`
	Description string        `json:"description"`
	Status      RequestStatus `json:"status"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// IsPending checks if request is awaiting approval
func (r *TrainingRequest) IsPending() bool {
	return r.Status == StatusPending
}

// IsApproved checks if request has been approved
func (r *TrainingRequest) IsApproved() bool {
	return r.Status == StatusApproved
}

// IsRejected checks if request has been rejected
func (r *TrainingRequest) IsRejected() bool {
	return r.Status == StatusRejected
}

// Approve changes status to approved
func (r *TrainingRequest) Approve() {
	r.Status = StatusApproved
	r.UpdatedAt = time.Now()
}

// Reject changes status to rejected
func (r *TrainingRequest) Reject() {
	r.Status = StatusRejected
	r.UpdatedAt = time.Now()
}
