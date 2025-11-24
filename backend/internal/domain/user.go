package domain

import "time"

// UserRole defines possible user roles in the system
type UserRole string

const (
	RoleEmployee UserRole = "employee"
	RoleAdmin    UserRole = "admin"
	RoleUser     UserRole = "user"
)

// User represents a system user (employee or administrator)
type User struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose in JSON
	Role         UserRole  `json:"role"`
	Department   *string   `json:"department,omitempty"`
	JobTitle     *string   `json:"jobTitle,omitempty"`
	Telegram     *string   `json:"telegram,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// IsAdmin checks if user has admin privileges
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsEmployee checks if user is a regular employee
func (u *User) IsEmployee() bool {
	return u.Role == RoleEmployee
}
