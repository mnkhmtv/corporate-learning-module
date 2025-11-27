package dto

import (
	"time"
)

// LearningRequestDTO represents embedded request info in learning response
type LearningRequestDTO struct {
	ID          string `json:"id"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
}

// LearningUserDTO represents embedded user info in learning response
type LearningUserDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LearningMentorDTO represents embedded mentor info in learning response
type LearningMentorDTO struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Telegram   *string `json:"telegram,omitempty"`
	JobTitle   string  `json:"jobTitle"`
	Experience *string `json:"experience,omitempty"`
}

// LearningProcessResponseDTO represents learning process response with embedded objects
type LearningProcessResponseDTO struct {
	ID        string                `json:"id"`
	Request   LearningRequestDTO    `json:"request"`
	User      LearningUserDTO       `json:"user"`
	Mentor    LearningMentorDTO     `json:"mentor"`
	Status    string                `json:"status"`
	StartDate time.Time             `json:"startDate"`
	EndDate   *time.Time            `json:"endDate,omitempty"`
	Plan      []LearningPlanItemDTO `json:"plan"`
	Feedback  *FeedbackDTO          `json:"feedback,omitempty"`
	Notes     *string               `json:"notes,omitempty"`
}
