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

// UpdatePlanDTO represents updating the entire learning plan
type UpdatePlanDTO struct {
	Plan []PlanItemDTO `json:"plan" binding:"required"`
}

// PlanItemDTO represents a plan item in request/response
type PlanItemDTO struct {
	ID        string `json:"id" binding:"required"`
	Text      string `json:"text" binding:"required"`
	Completed bool   `json:"completed"`
}

// UpdateNotesDTO represents updating learning notes
type UpdateNotesDTO struct {
	Notes string `json:"notes"`
}

// CompleteLearningDTO represents completing learning with feedback
type CompleteLearningDTO struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required"`
}
