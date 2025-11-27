package dto

// CreateRequestDTO represents training request creation input
type CreateRequestDTO struct {
	Topic       string `json:"topic" binding:"required"`
	Description string `json:"description" binding:"required"`
}
