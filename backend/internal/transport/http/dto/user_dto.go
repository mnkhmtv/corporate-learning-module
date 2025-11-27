package dto

// UpdateUserDTO represents user update request
type UpdateUserDTO struct {
	Name       *string `json:"name"`
	Email      *string `json:"email"`
	Password   *string `json:"password"`
	Department *string `json:"department"`
	JobTitle   *string `json:"jobTitle"`
	Telegram   *string `json:"telegram"`
}
