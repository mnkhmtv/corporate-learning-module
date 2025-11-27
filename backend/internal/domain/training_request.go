package domain

import (
	"errors"
	"time"
)

// RequestStatus represents the status of a training request
type RequestStatus string

const (
	RequestPending  RequestStatus = "pending"
	RequestApproved RequestStatus = "approved"
	RequestRejected RequestStatus = "rejected"
)

// TrainingRequest represents a request for training or mentorship
type TrainingRequest struct {
	ID          string        `json:"id"`
	UserID      string        `json:"userId"`
	Topic       string        `json:"topic"`
	Description string        `json:"description"`
	Status      RequestStatus `json:"status"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`

	// ↓ Новые поля из JOIN с users (для response DTO)
	UserName     string  `json:"-"` // Не показывать в JSON напрямую
	UserJobTitle *string `json:"-"`
	UserTelegram *string `json:"-"`
	// ↑
}

// Validate checks if the training request is valid
func (tr *TrainingRequest) Validate() error {
	if tr.UserID == "" {
		return errors.New("user ID is required")
	}
	if tr.Topic == "" {
		return errors.New("topic is required")
	}
	if tr.Description == "" {
		return errors.New("description is required")
	}
	return nil
}

// IsPending checks if the request is pending
func (tr *TrainingRequest) IsPending() bool {
	return tr.Status == RequestPending
}

// IsApproved checks if the request is approved
func (tr *TrainingRequest) IsApproved() bool {
	return tr.Status == RequestApproved
}

// IsRejected checks if the request is rejected
func (tr *TrainingRequest) IsRejected() bool {
	return tr.Status == RequestRejected
}

// Approve marks the request as approved
func (tr *TrainingRequest) Approve() {
	tr.Status = RequestApproved
}

// Reject marks the request as rejected
func (tr *TrainingRequest) Reject() {
	tr.Status = RequestRejected
}
