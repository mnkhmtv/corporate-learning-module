package domain

import "time"

// Mentor represents a training mentor in the system
type Mentor struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	JobTitle   string    `json:"jobTitle"`
	Experience *string   `json:"experience,omitempty"`
	Workload   int       `json:"workload"` // 0-5 scale
	Email      string    `json:"email"`
	Telegram   *string   `json:"telegram,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// IsAvailable checks if mentor has capacity for new students
func (m *Mentor) IsAvailable() bool {
	return m.Workload < 5
}

// CanTakeStudent checks if mentor can accept one more student
func (m *Mentor) CanTakeStudent() bool {
	return m.Workload <= 4
}

// IncrementWorkload increases mentor's workload by 1
func (m *Mentor) IncrementWorkload() {
	if m.Workload < 5 {
		m.Workload++
	}
}

// DecrementWorkload decreases mentor's workload by 1
func (m *Mentor) DecrementWorkload() {
	if m.Workload > 0 {
		m.Workload--
	}
}
