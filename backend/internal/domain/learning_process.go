package domain

import (
	"errors"
	"time"
)

// LearningStatus defines the status of a learning process
type LearningStatus string

const (
	LearningActive    LearningStatus = "active"
	LearningCompleted LearningStatus = "completed"
)

// Feedback represents feedback for completed learning
type Feedback struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

// LearningProcess represents an active or completed learning session
type LearningProcess struct {
	ID        string             `json:"id"`
	RequestID string             `json:"requestId"`
	UserID    string             `json:"userId"`
	MentorID  string             `json:"mentorId"`
	Status    LearningStatus     `json:"status"`
	StartDate time.Time          `json:"startDate"`
	EndDate   *time.Time         `json:"endDate,omitempty"`
	Plan      []LearningPlanItem `json:"plan"`
	Feedback  *Feedback          `json:"feedback,omitempty"`
	Notes     *string            `json:"notes,omitempty"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`

	RequestTopic       string  `json:"-"`
	RequestDescription string  `json:"-"`
	UserName           string  `json:"-"`
	MentorName         string  `json:"-"`
	MentorTelegram     *string `json:"-"`
	MentorJobTitle     string  `json:"-"`
	MentorExperience   *string `json:"-"`
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
	return nil
}

// UpdatePlanItem updates an existing plan item by ID
func (lp *LearningProcess) UpdatePlanItem(id string, text string, completed bool) error {
	for i := range lp.Plan {
		if lp.Plan[i].ID == id {
			if text != "" {
				lp.Plan[i].Text = text
			}
			lp.Plan[i].Completed = completed
			return nil
		}
	}
	return ErrPlanItemNotFound
}

// TogglePlanItem toggles completion status of a plan item
func (lp *LearningProcess) TogglePlanItem(id string) error {
	for i := range lp.Plan {
		if lp.Plan[i].ID == id {
			lp.Plan[i].Toggle()
			return nil
		}
	}
	return ErrPlanItemNotFound
}

// RemovePlanItem removes an item from the plan by ID
func (lp *LearningProcess) RemovePlanItem(id string) error {
	for i, item := range lp.Plan {
		if item.ID == id {
			lp.Plan = append(lp.Plan[:i], lp.Plan[i+1:]...)
			return nil
		}
	}
	return ErrPlanItemNotFound
}

// GetPlanItem returns a specific plan item by ID
func (lp *LearningProcess) GetPlanItem(id string) (*LearningPlanItem, error) {
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
	lp.Feedback = &Feedback{
		Rating:  rating,
		Comment: comment,
	}
	now := time.Now()
	lp.EndDate = &now

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

// Validate feedback
func (f *Feedback) Validate() error {
	if f.Rating < 1 || f.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	if f.Comment == "" {
		return errors.New("comment cannot be empty")
	}
	return nil
}
