package dto

// RegisterDTO represents registration request
type RegisterDTO struct {
	Name       string  `json:"name" binding:"required"`
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required,min=6"`
	Department *string `json:"department"`
	JobTitle   *string `json:"jobTitle"`
	Telegram   *string `json:"telegram"`
}

// LoginDTO represents login request
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login output
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
