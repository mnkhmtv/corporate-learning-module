package dto

// AddPlanItemDTO represents adding an item to learning plan
type AddPlanItemDTO struct {
	Text string `json:"text" binding:"required"`
}

// UpdatePlanItemDTO represents updating a plan item
type UpdatePlanItemDTO struct {
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// CompleteLearningDTO represents completing learning with feedback
type CompleteLearningDTO struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required"`
}
