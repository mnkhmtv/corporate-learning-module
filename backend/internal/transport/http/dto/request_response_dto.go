package dto

import (
	"time"
)

// RequestUserDTO represents embedded user info in request response
type RequestUserDTO struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	JobTitle *string `json:"jobTitle,omitempty"`
	Telegram *string `json:"telegram,omitempty"`
}

// TrainingRequestResponseDTO represents request response with embedded user
type TrainingRequestResponseDTO struct {
	ID          string         `json:"id"`
	User        RequestUserDTO `json:"user"`
	Topic       string         `json:"topic"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}
