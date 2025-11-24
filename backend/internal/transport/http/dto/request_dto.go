package dto

// CreateRequestDTO represents training request creation input
type CreateRequestDTO struct {
	Topic       string `json:"topic" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// AssignMentorDTO represents mentor assignment input
type AssignMentorDTO struct {
	MentorID string `json:"mentorId" binding:"required"`
}
