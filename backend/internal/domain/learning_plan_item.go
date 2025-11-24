package domain

import "errors"

// LearningPlanItem represents a single task or milestone in a learning plan
type LearningPlanItem struct {
	ID        uint8  `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// NewLearningPlanItem creates a new plan item with specified ID
func NewLearningPlanItem(id uint8, text string) (*LearningPlanItem, error) {
	if text == "" {
		return nil, errors.New("plan item text cannot be empty")
	}

	return &LearningPlanItem{
		ID:        id,
		Text:      text,
		Completed: false,
	}, nil
}

// MarkCompleted marks the item as completed
func (item *LearningPlanItem) MarkCompleted() {
	item.Completed = true
}

// MarkIncomplete marks the item as incomplete
func (item *LearningPlanItem) MarkIncomplete() {
	item.Completed = false
}

// Toggle switches the completion status
func (item *LearningPlanItem) Toggle() {
	item.Completed = !item.Completed
}

// Validate checks if the plan item is valid
func (item *LearningPlanItem) Validate() error {
	if item.Text == "" {
		return errors.New("plan item text is required")
	}
	return nil
}
