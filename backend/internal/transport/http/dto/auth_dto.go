package dto

// RegisterRequest represents registration input
type RegisterRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	Role       string `json:"role"`
	Department string `json:"department"`
	JobTitle   string `json:"jobTitle"`
	Telegram   string `json:"telegram"`
}

// LoginRequest represents login input
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login output
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
