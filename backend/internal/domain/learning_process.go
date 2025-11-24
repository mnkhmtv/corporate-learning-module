package domain

import "time"

// LearningStatus defines the status of a learning process
type LearningStatus string

const (
	LearningActive    LearningStatus = "active"
	LearningCompleted LearningStatus = "completed"
)

// LearningProcess represents an active or completed learning session
type LearningProcess struct {
	ID              uint64             `json:"id"`
	RequestID       uint64             `json:"requestId"`
	UserID          uint64             `json:"userId"`
	MentorID        uint64             `json:"mentorId"`
	Status          LearningStatus     `json:"status"`
	Plan            []LearningPlanItem `json:"plan"`
	Notes           *string            `json:"notes,omitempty"`
	FeedbackRating  *int               `json:"feedbackRating,omitempty"`
	FeedbackComment *string            `json:"feedbackComment,omitempty"`
	CreatedAt       time.Time          `json:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt"`
	CompletedAt     *time.Time         `json:"completedAt,omitempty"`
}

// IsActive checks if learning process is ongoing
func (lp *LearningProcess) IsActive() bool {
	return lp.Status == LearningActive
}

// IsCompleted checks if learning process is finished
func (lp *LearningProcess) IsCompleted() bool {
	return lp.Status == LearningCompleted
}

// AddPlanItem adds a new item to the learning plan
func (lp *LearningProcess) AddPlanItem(item LearningPlanItem) error {
	if err := item.Validate(); err != nil {
		return err
	}
	lp.Plan = append(lp.Plan, item)
	lp.UpdatedAt = time.Now()
	return nil
}

// GetNextPlanItemID returns the next available ID for a new plan item
func (lp *LearningProcess) GetNextPlanItemID() (uint8, error) {
	if len(lp.Plan) >= 255 {
		return 0, ErrPlanLimitReached
	}

	maxID := uint8(0)
	for _, item := range lp.Plan {
		if item.ID > maxID {
			maxID = item.ID
		}
	}

	return maxID + 1, nil
}

// UpdatePlanItem updates an existing plan item by ID
func (lp *LearningProcess) UpdatePlanItem(id uint8, text string, completed bool) error {
	for i := range lp.Plan {
		if lp.Plan[i].ID == id {
			if text != "" {
				lp.Plan[i].Text = text
			}
			lp.Plan[i].Completed = completed
			lp.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrPlanItemNotFound
}

// TogglePlanItem toggles completion status of a plan item
func (lp *LearningProcess) TogglePlanItem(id uint8) error {
	for i := range lp.Plan {
		if lp.Plan[i].ID == id {
			lp.Plan[i].Toggle()
			lp.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrPlanItemNotFound
}

// RemovePlanItem removes an item from the plan by ID
func (lp *LearningProcess) RemovePlanItem(id uint8) error {
	for i, item := range lp.Plan {
		if item.ID == id {
			lp.Plan = append(lp.Plan[:i], lp.Plan[i+1:]...)
			lp.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrPlanItemNotFound
}

// GetPlanItem returns a specific plan item by ID
func (lp *LearningProcess) GetPlanItem(id uint8) (*LearningPlanItem, error) {
	for i := range lp.Plan {
		if lp.Plan[i].ID == id {
			return &lp.Plan[i], nil
		}
	}
	return nil, ErrPlanItemNotFound
}

// Complete marks the learning process as completed with feedback
func (lp *LearningProcess) Complete(rating int, comment string) error {
	if !lp.IsActive() {
		return ErrLearningNotActive
	}

	if rating < 1 || rating > 5 {
		return ErrInvalidRating
	}

	lp.Status = LearningCompleted
	lp.FeedbackRating = &rating
	lp.FeedbackComment = &comment
	now := time.Now()
	lp.CompletedAt = &now
	lp.UpdatedAt = now

	return nil
}

// GetProgress returns the percentage of completed plan items
func (lp *LearningProcess) GetProgress() float64 {
	if len(lp.Plan) == 0 {
		return 0.0
	}

	completed := 0
	for _, item := range lp.Plan {
		if item.Completed {
			completed++
		}
	}

	return float64(completed) / float64(len(lp.Plan)) * 100
}

// GetCompletedItemsCount returns the number of completed items
func (lp *LearningProcess) GetCompletedItemsCount() int {
	count := 0
	for _, item := range lp.Plan {
		if item.Completed {
			count++
		}
	}
	return count
}
